package utils

import (
	"testing"
)

const (
	raw = `Milvus 是全球最快的向量数据库，在最新发布的Milvus 2.2 benchmark中，Milvus相比之前的版本，取得了50%以上的性能提升。值得一提的是，在Master branch的最新分支中，Milvus的性能又更进一步，在1M向量串行执行的场景下取得了3ms以下的延迟，整体QPS甚至超过了ElasticSearch的10倍。那么，如何使用Milvus才能达到理想的性能呢？本文暂且不提社区大神贡献的黑科技优化，先聊聊使用Milvus过程中一些经验，以及如何进行性能调优。经验1 - 合理的预计数据量，表数目大小，QPS参数等指标在部署Milvus之前，首先需要决定机器的资源，规格，以及一些依赖的资源，以下是你需要考虑的因素：有多少张表？每张表的数据量有多少？每张表的QPS需求有多少？是否需要存标量字段，如果有字符串，字符串的平均长度？是否有删除和流式插入，每天大概有多少比例的数据需要倍更新？基于以上因素，有以下经验结论可以遵循：节点资源占用可以通过sizing tool 进行计算，通常情况下下8G内存可以支持超过5m的128dim向量数据和1m的768dim数据。默认情况下，Milvus会创建256个消息队列topic。如果表数目比较少，可以调整rootCoord.dmlChannelNum减少topic数目降低消息队列负载。默认情况下，每个collection会使用2个消息队列topic（shard），如果写入非常大或者数据量极大，需要调整collection的shard数目。建议每个shard写入/删除不超过10M/s，单个shard的数据量不大于1B向量， shard数目过大也会影响写入性能，因此不建议单表超过8个shard。根据benchmark结果计算需要的CPU资源。对于小数据量场景（小于5m），使用多副本可以扩展查询性能，但建议副本数目不要超过10个。对于中大数据量场景，通常扩容querynode就可以自动负载均衡，不需要使用多副本提升QPS.所有的标量字段目前也会加载进内存中，也会消耗内存，请在容量规划时预留原始数据类型两倍以上的内存。Milvus 在存储数据的过程中，存在较多冗余数据（https://github.com/milvus-io/milvus/issues/20453）. 考虑到Minio的2，4纠删码存在两副本冗余，我们建议Minio至少包含6倍以上的数据的磁盘存储。同时Pulsar/Kafka需要包含近五天写入量三倍的存储。合理调整数据的保留时间和GC时间可以很大程度上减少磁盘的使用，默认情况下数据会被保留5天。个人建议适当缩短数据过期时间，但尽可能保留1天以上以避免数据丢失或误删除。Etcd作为Milvus的元信息存储和服务发现节点，请尽可能使用ssd磁盘并独立部署。通常Etcd的内存使用不会超过4GB，通过调整参数可以较快的清理etcd中的历史版本减少内存使用。Pulsar/Kafka作为Milvus的日志存储，其依赖的zookeeper集群对性能要求也比较高，建议使用SSD并独立部署。经验2 - 选择合适的索引类型和参数索引的选择对于向量召回的性能至关重要，Milvus支持了Annoy，Faiss，HNSW，DiskANN等多种不同的索引，用户可以根据对延迟，内存使用和召回率的需求进行选择。索引的选择步骤一般如下：1） 是否需要精确结果？只有Faiss的Flat索引支持精确结果，但需要注意Flat索引检索速度很慢，查询性能通常比其他Milvus支持的索引类型低两个数量级以上，因此只适合千万级数据量的小查询（Flat on GPU 已经在路上了，敬请期待）2）数据量是否能加载进内存？对于大数据量，内存不足的场景，Milvus提供两种解决方案：DiskANNDiskANN依赖高性能的磁盘索引，借助NVMe磁盘缓存全量数据，在内存中只存储了量化后的数据.DiskANN适用于对于查询Recall要求较高，QPS不高的场景DiskANN示意图 DiskANN的关键参数：       search_list:  search_list越大，recall越高而性能越差。search_list的大小不应该小于K。 而对于较小的K，推荐把search_list和K的比值设置得相对大一些, 这个比值随着K增大可以逐渐靠近1。  2. IVF_PQ对于精确度要求不高的场景或者性能要求极高的场景IVF PQ的核心是两个算法，IVF + PQ量化，其中量化可以大幅减少向量的占用内存量IVF参数nlist：一般建议nlist =4*sqrt(N),对于Milvus而言，一个Segment默认是512M数据，对于128dim向量而言，一个segment包含100w数据，因此最佳nlist在1000左右。nprobe：nprobe可以Search时调整搜索的数据量，nprobe越大，recall越高，但性能越差。具体的nprobe需要根据查询的精度要求决定，从nprobe=16开始会是一个不错的尝试。 PQ参数M:  向量做PQ的分段数目，一般建议设置为向量维数的1/4，M取值越小内存占用越小，查询速度越快，精度也变得更加低。Nbits: 每段量化器占用的bit数目，默认为8，不建议调整3） 构建索引和内存资源是否充足性能优先，选择HNSW索引HNSW索引是目前Milvus支持的性能最快的索引，我们的测试报告也是基于HNSW作为测试依据HNSW内存的开销较高，通常需要原始向量的1.5 - 2倍以上内存HNSW参数M：表示在建表期间每个向量的边数目，M越大，内存消耗越高，在高维度的数据集下查询性能会越好。通常建议设置在8-32之间ef_construction： 控制索引时间和索引准确度，ef_construction越大构建索引越长，但查询精度越高。要注意ef_construction 提高并不能无限增加索引的质量，常见的ef_construction参数为128.ef: 控制搜索精确度和搜索性能，注意ef必须大于K。资源优先，选择IVF_FLAT或者IVF_SQ8索引IVF索引在Milvus分片之后也能拿到比较不错的召回率，其内存占用和建索引速度相比HNSW都要低很多IVF_SQ8相比IVF，将向量数据从float32转换为了int8，可以减少4倍的内存用量，但对召回率有较大影响，如果要求95%以上的召回精度不建议使用。IVF类索引的参数跟IVFPQ类似，这里就不做过多的介绍了。检索时，Milvus的查询一致性也会对查询造成较大影响。通常情况，对于一致性要求较高的场景，建议使用最终一致性或者有界一致，默认情况下Milvus选择有界一致性，窗口为3s。`
)

func TestSplit(t *testing.T) {
	result, err := TextSplit(raw)
	if err != nil {
		t.Fatal(err)
	}
	for _, v := range result {
		t.Logf("content:%v, length=%v, tokens=%v\n", v.Content, v.ContentLength, v.ContentTokens)
	}
}
