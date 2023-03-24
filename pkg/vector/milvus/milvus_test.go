package milvus

import (
	"context"
	"fmt"
	"testing"

	"github.com/zach030/fable/pkg/vector/model"

	"github.com/milvus-io/milvus-sdk-go/v2/entity"

	"github.com/stretchr/testify/assert"
)

var (
	milvusClient *MilvusClient
	addr         = ""
	user         = ""
	pwd          = ""
	inputVector  = []float32{0.016277889, -0.007834382, -0.010743925, -0.013763629, -0.014139472, 0.03146713, -0.011651132, -0.014865237, -0.004066232, -0.016122367, 0.014087631, 0.00756222, -0.004691557, 0.0016548431, -0.004668877, -0.0018743548, 0.024144672, 0.012493539, 0.013763629, 0.021811852, -0.005446483, 0.022615379, 0.00504472, 0.0080547035, -0.0038783108, -0.004516596, 0.0025644803, -0.03154489, 0.020697284, -0.005812606, 0.02012704, -0.018999511, -0.0132711455, -0.006667973, -0.020865764, -0.006687413, -0.012227857, 0.004007912, -0.0026033607, -0.030559922, -0.00057550956, 0.0051192404, -0.004548996, -0.007102136, -0.036106847, 0.01128177, 0.020476962, -0.0038847907, -0.006165769, -0.0063731303, -0.0048956787, 0.020878725, -0.014100592, -0.02789014, -0.018442227, 0.009661756, -0.017910862, 0.028045662, 0.011508571, -0.0027410616, -0.015487323, -0.016135328, 0.013932111, 0.024300192, -0.007685341, 0.016731493, -0.028253024, 0.012098256, 0.0015827526, -0.0027329617, 0.039450552, -0.0018419545, 0.01645933, -0.011521531, 0.0008634668, -0.0012757601, -0.008346306, 0.0011631692, 0.0014312813, -0.010335682, 0.010627284, -0.0067003733, -0.009026712, 0.030508082, 0.011910334, 0.0035024676, 0.017042534, 0.009130392, -0.014878198, -0.0067327735, -0.0021967373, 0.012739781, 0.02913431, 0.023198584, -0.02649045, -0.011929775, 0.013724749, 0.013348905, 0.0026811212, -0.01132713, 0.010659684, 0.010802246, -0.029885996, -0.01006352, -0.020697284, -0.01511148, 0.009454395, 0.0086638285, 0.002906303, -0.014230193, -0.019958558, 0.030922804, 0.0061333687, -0.024390914, 0.0012279698, -0.012895302, -0.024831556, -0.014839318, -0.008223185, -0.02269314, 0.01885695, 0.01901247, 0.023017142, -0.01259074, 0.012195457, 0.003956071, -0.021306409, -0.024831556, 0.009499756, -0.020489922, 0.044738274, 0.010815206, 0.0024073392, -0.0027345817, -0.009298874, 0.017081415, -0.036029086, 0.03133753, -0.017366538, -0.0048697586, 0.0014758317, 0.029885996, -0.008113024, -0.007063256, -0.004247674, 0.031078326, 0.0022906982, -0.011735373, 0.029600874, -0.025699884, 0.027864221, -0.007659421, 0.003149305, 0.002164337, 0.012493539, 0.021941453, 0.02646453, 0.006687413, -0.013569227, 0.022187695, 0.007944543, -0.0025580002, 0.017418377, -0.004088912, 0.037506536, 0.021163847, 0.024261313, -0.010147761, -0.00629213, -0.003667709, -0.01883103, 0.0005544494, -0.003528388, 0.029160231, -0.015539163, 0.008599028, 0.021150887, -0.006687413, -0.024131712, -0.015850205, -0.04942983, -0.007659421, 0.007458539, 0.0070956564, 0.0125389, -0.0064055305, 0.04966311, 0.004412915, 0.004173153, -0.012286177, 0.011152169, 0.013374826, -0.012752741, -0.002255058, -0.66646034, -0.026645971, -0.013970991, -0.00884527, 0.019193912, 0.010232001, 0.015318842, 0.0064930115, -0.002911163, -0.009473835, 0.007341898, -0.003664469, -0.009804318, -0.010296801, 0.0058806464, -0.02014, 0.0092729535, 0.0036255887, -0.0035931885, 0.023833629, -0.012072336, 0.006139849, -0.018014543, -0.0073030177, -0.002161097, 0.0021011566, -0.010970727, -0.0038847907, -0.019232793, 0.015681725, -0.03177817, 0.028693667, 0.009674717, 0.0050933203, 0.0352774, -0.01631677, -0.01258426, 0.027864221, 0.0015754625, 0.05038888, -0.03193369, 0.01651117, -0.0058547263, 0.010737445, -0.01653709, 0.019362394, 0.015772445, -0.012409299, 0.018986551, -0.016161248, 0.03172633, 0.006551332, 0.009020232, 0.014087631, -0.0002780347, 0.019180952, 0.012474099, 0.00035437781, 0.018234864, -0.0017026335, -0.020593602, 0.013698828, -0.014917078, -0.018766228, -0.005692725, 0.009130392, -0.014968919, -0.0021594772, 0.0078019816, -0.011955694, 0.011612252, 0.0036353087, -0.005329842, 0.003290246, 0.01127529, -0.0051192404, 0.04574916, -0.019336473, 0.027475417, 0.00758814, -0.0034927477, -0.0016362129, -0.017314697, -0.0024381194, 0.03504412, -0.017923823, -0.015798366, -0.008061184, 0.014657876, -0.00043699847, 0.002781562, 0.0015722224, 0.0039625512, -0.0063958107, 0.01649821, 0.042042572, -0.009480315, -0.0016848133, -0.00066542026, 0.013543307, -0.0045749163, -0.00044023848, 0.002044456, 0.00045562861, 0.031389367, -0.005346042, -0.008981351, 0.0032011454, 0.03307418, -0.0302748, 0.0028317824, -0.01505964, -0.044971555, -0.005378443, 0.0034441473, -0.023107862, 0.007717741, -0.0047336775, 0.017301736, -0.012908262, 0.015565083, -0.035095956, 0.018416306, 0.00079340127, 0.016368609, 0.012500019, -0.018532947, -0.024455713, -0.014722677, -0.013245225, 0.015889086, 0.0010935087, 0.015344761, -0.00021910673, 0.024287233, 0.01385435, 0.008132464, -0.022291377, 0.005942207, -0.027449498, -0.001901895, 0.005346042, -0.020658404, 0.007439099, -0.0023652187, -0.031622652, -0.034266513, -0.034862675, -0.0022485778, 0.015305881, 0.0047304374, -0.0020476962, -0.027345816, 0.0015066119, -0.016718533, -7.7659366e-05, -0.02400211, -0.00885175, -0.013335946, -0.005478883, -0.006979015, -0.0012798101, -0.0134007465, -0.011605772, -0.012882342, -0.030534001, -0.029419433, 0.034266513, -0.0029289832, -0.021734092, -0.010322722, 0.011528011, -0.0014086012, 0.009182232, -0.017949741, -0.0067327735, -0.02036032, -0.023600347, -0.0013729609, 0.005164601, 0.026321968, -0.0017204536, -0.0115733715, -0.00634721, 0.020489922, 0.019310553, -0.00027965472, 0.013102664, 0.00030820744, 0.008579588, -0.014541235, 0.00378759, 0.017301736, -0.00035741532, 0.00506416, 0.007918622, -0.014411634, 0.012435218, 0.017301736, -0.0047563575, 0.035147797, -0.003170365, 0.021967374, -0.009674717, 0.0063731303, -0.005161361, 0.019712316, -0.01630381, -0.018442227, 0.028045662, -0.010828165, -0.018053424, -0.018286705, -0.01752206, 0.0052812416, 0.031104246, -0.026360849, -0.0014839318, -0.017003655, -0.0036515088, -0.017198056, 0.0074002184, 0.0033793468, -0.015655804, -0.004950759, 0.010828165, -0.014981879, 0.023276344, 0.008449987, -0.05165897, -0.0016993935, 0.008501828, -0.01885695, 0.015746525, -0.008022304, 0.0032691858, -0.004782278, -0.0136599485, 0.033437066, -0.011229929, -0.008326866, 0.009895039, 0.031259768, 0.0029565233, -0.004169913, 0.01394507, 0.02147489, 0.0097524775, -0.03525148, -0.006165769, -0.006933655, -0.01759982, 0.0054983236, -0.004011152, 0.023768827, -0.0091433525, 0.024727875, -0.0008334966, 0.007860302, 0.017859021, 0.005423803, 0.0004799288, -0.00046332367, -0.007477979, -0.0032173456, -0.009920958, -0.016187169, 0.00884527, -0.0086638285, -0.0016702332, -0.028382625, -0.00628889, 0.031804092, -0.0054853633, -0.0035964285, -0.0003033474, 0.008190785, 0.019634556, 0.01137897, 0.018105263, -0.016368609, 0.0039755115, -0.0071474966, 0.016187169, -0.0052261613, -0.029445354, -0.019362394, 0.01640749, 0.0120528955, 0.0067262934, 0.0012943903, 0.0064962516, -0.0019683156, -0.0304044, -0.011832573, -0.027060695, 0.005537204, -0.030508082, -0.000501394, 0.01772942, 0.002162717, 0.0043707946, -0.018766228, -0.022991221, 0.025803564, 0.025142599, -0.002771842, -0.0024235393, -0.001645933, -0.022291377, -0.041705612, -0.01893471, 0.0060102474, 0.009434954, 0.0018208944, -0.0019310553, 0.01511148, 0.0026130807, 0.034862675, 0.026438609, 0.004429115, -0.024209471, 0.010031119, -0.008605508, 0.08771397, 0.036184605, 0.010983687, 0.0020671363, -0.022433938, 0.009149833, -0.008087104, -0.009000791, -0.009694157, 0.0028755227, 0.006577252, 0.0069984556, -0.008074144, 0.0006852654, 0.034810837, 0.01511148, 0.015331801, -0.006531892, 0.014942998, -0.016783332, -0.019673435, 0.0072511774, -0.0035219078, 0.010264401, 0.02278386, -0.00379083, 0.0032465057, 0.014100592, 0.00632777, -0.0075946203, -0.008002863, 0.038387824, -0.0036515088, -0.00379407, 0.0033501866, 0.019479034, -0.019751197, 0.002653581, 0.006421731, -0.013199865, 0.024339072, 0.0018063143, 0.012765701, -0.018416306, 0.011100328, -0.017042534, -0.017327657, 0.018740308, -0.039968956, -0.016601892, 0.031259768, 0.016563011, -0.03390363, -0.0010975586, -0.00505444, 0.046137966, 0.0010473382, 0.010568963, -0.00254504, -0.0057575256, -0.020904645, -0.014282033, 0.0065254117, 0.00023206684, -0.011566891, -0.018260784, 4.796757e-06, -0.022835702, -0.015733564, 0.015007799, -0.014593076, -0.0328409, -0.013750669, 0.011003127, 0.032529857, 0.0019229553, -0.0023587386, -0.02396323, 0.016718533, -0.006136609, 0.00505444, -0.0055242437, 0.008119504, -0.020010399, 0.010452323, -0.021889614, 0.015681725, -0.028019741, -0.0052358815, 0.01762574, 0.01007648, -0.0008351166, 0.033825867, -0.018701429, 0.0020622762, -0.012033455, -0.004247674, 0.017418377, -0.004672117, 0.00876751, 0.043001622, 0.00032461007, 0.00031407998, -0.021034246, 0.01903839, 0.0013267905, 0.025440682, 0.026749652, -0.00887767, 0.0027443017, -0.0137895495, -0.021008326, 0.0033096862, 0.0062014093, 0.016057568, -0.0043545943, -0.010031119, -0.005929247, -0.019258713, -0.0056344047, 0.0019229553, -0.037558377, 0.014865237, 0.0015819425, 0.00632129, 0.008676789, 0.009525675, -0.0030051237, -0.021112008, 0.0055242437, -0.04274242, 0.00029828487, -0.0060102474, -0.015046679, -0.012117696, -0.0045975964, 0.008994311, 0.0023182384, -0.023755869, -0.015798366, -0.019686395, 0.012150096, -0.010283842, -0.034888595, -0.0152022, -0.03675485, -0.0009914478, -0.013258185, -0.0022356177, 0.034422033, -0.03307418, -0.007076216, -0.005274762, -0.005692725, -0.01639453, -0.028304864, -0.009927439, -0.036858533, 0.055106357, 0.026957013, 0.032270655, -0.012279698, 0.027345816, -0.004309234, -0.013070264, -0.01769054, 0.004020872, 0.011955694, -0.020736163, 0.0047531174, 0.016757412, -0.0057380856, 0.021721132, -0.005427043, 0.025557322, -0.00503176, -0.010854086, -0.006062088, -0.022084014, -0.04201665, -0.02388547, -0.01902543, -0.002149757, 0.010212561, -0.018208943, -0.003421467, 0.012422258, 0.022226576, 0.005436763, 0.007789022, 0.017547978, 0.010581924, 0.03636605, 0.004429115, -0.0022080776, -0.008890631, -0.010718005, -0.0036417888, -0.017107336, 0.024377953, 0.0042606336, -0.00096795766, 0.034266513, -0.018170064, 0.010983687, 0.032063294, -0.0011858494, -0.0048535583, -0.009104472, -0.025985006, -0.018584788, -0.029834157, -0.03781758, -0.035795804, 0.01385435, 0.02016592, -0.020554723, 0.019103192, -0.007348378, -0.02287458, -0.003544588, -0.0010368082, 0.031156087, -0.0008334966, 0.012214896, 0.024274273, -0.007322458, -0.019090232, 0.0135303475, -0.020826885, -0.010718005, 0.019116152, 1.1580561e-06, 0.00077436614, 0.008229665, 0.0042379536, 0.01138545, -0.028408544, -0.02538884, 0.015228121, -0.006207889, 0.020982407, -0.012441698, -0.007063256, -0.013983951, -0.0065092114, -0.02265426, 0.008624949, 0.006535132, -0.020010399, -0.006211129, 0.030093359, -0.011560411, 0.0019683156, 0.0032967262, 0.011988095, -0.002408959, 0.00079259125, 0.019336473, -0.0040241117, -0.004924839, 0.02669781, -0.009026712, -0.0053687226, 0.009726557, 0.014035791, -0.008028784, -0.0072706174, -0.014320914, 0.025414761, -0.013569227, -0.009026712, 0.0077436613, 0.0007200957, 0.013556267, 0.0035186678, 0.0021027767, -0.016899973, -0.008145425, -0.00376491, 0.021941453, 0.0034668276, -0.008560148, 0.0041148327, -0.012150096, -0.010815206, -0.0020671363, -0.009875598, 0.027553178, -0.0328409, -0.00072738575, -0.0097071165, 0.009454395, 0.016070528, 0.010724485, 0.0021934973, 0.0029743435, 0.024339072, -0.0043707946, 0.010594884, 0.024507554, 0.033774026, -0.0021886374, 0.018364465, -0.009357194, -0.03543292, 0.016161248, -0.008449987, -0.020412162, -0.009318314, -0.00056538446, -0.0031671252, -0.008961911, 0.027034774, 0.005197001, 0.016174208, -0.0029532835, -0.0063860905, -0.03175225, 0.0019083751, -0.010517123, -0.0015989527, 0.013374826, -0.0013470406, -0.0024624194, -0.025505481, -0.010983687, 0.00063950004, -0.047459897, -0.014282033, -0.006172249, 0.034966357, 0.0012741401, -0.009052631, -0.012338018, -0.00095823756, -0.024403874, -0.0014377614, -0.009117432, 0.014424594, 0.0024348793, 0.006447651, -0.0013316505, 0.027993822, -0.0014369513, -0.024663076, -0.019103192, -0.022239536, -0.012305617, -0.020839846, -0.01759982, 0.0068558943, 0.040150397, -0.0020833365, -0.002271258, 0.004276834, -0.0075946203, -0.009745997, -0.0015924727, 0.005167841, 0.024870437, 0.012279698, 0.00023793938, 0.016822213, -0.003411747, -0.002653581, -0.011696492, -0.0072770976, -7.669748e-06, 0.00253046, 0.032400258, -0.020839846, 0.0014952718, -0.002033116, 0.025881326, 0.0011169988, 0.007873262, 0.0029678636, -0.018338544, 0.02669781, -0.0030715442, 0.0018403346, 0.0066485326, 0.018066384, 0.00091773726, -0.03265946, -0.0019359153, 0.0046073166, 0.009473835, -0.009499756, 0.0073872586, 0.025622122, 0.009778397, -0.01629085, -0.016938854, -0.006842934, 0.0061042085, -0.014139472, 0.017262857, 0.004092152, -0.0012798101, 0.00013689109, 0.009039671, -0.00046980372, 0.020632483, 0.0064994916, -0.0017253137, -0.018170064, -0.012519459, -0.02409283, 0.052410655, -0.017871981, 0.013219304, -0.027086614, -0.0121824965, -0.044971555, -0.018221904, -0.009577516, -0.01653709, 0.01753502, -0.008611988, 0.02789014, -0.007309498, 0.01008944, 0.012512979, -0.0046008364, 0.0056506046, 0.005942207, 0.03022296, 0.0061268886, 0.017327657, 0.00078084617, 0.01621309, -0.009007271, 0.01011536, 0.013141544, 0.25111493, -0.019945597, -0.015889086, 0.035536602, -0.0011655992, -0.0009914478, 0.0076983008, -0.0051095206, -0.0012846702, 0.022019215, -0.010517123, 0.0126166595, -0.01649821, -0.0054659233, 0.0203344, 0.016796293, -0.033488907, 0.0032967262, -0.021591531, -0.0018581547, 0.019103192, -0.008456467, -0.029730475, -0.0038912708, 0.01766462, -0.0053428025, 0.0042282334, 0.0203344, 0.021254567, 0.015681725, 0.0021853973, -0.0039463514, 0.0029532835, 0.0019391554, -0.008028784, 0.005922767, 0.00379407, 0.00016210253, 0.008041744, 0.033696268, 0.003155785, 0.016226048, -0.011191049, -0.0005200241, 0.016096447, -0.015772445, -0.008035264, -0.00037179294, 0.012564819, 0.016990695, -0.0152022, 0.012286177, 0.039943036, 0.018364465, -0.001901895, -0.026257169, -0.008547188, 0.022498738, 0.0032270656, 0.012901782, -0.026049806, 0.016575972, -0.020671364, 0.01516332, -0.0076399804, -0.0045068758, -0.034733076, 0.03672893, 0.035536602, -0.00879991, -0.01007648, -0.0028042423, -0.010536564, 0.008268545, -0.0076140603, -0.025544362, 0.025764683, 0.014424594, 0.0070956564, 0.03779166, -0.0038912708, -0.019453114, 0.0063958107, -0.02147489, 0.011819613, -0.051399767, 0.026982933, -0.008287986, -0.0106208045, 0.008411107, -0.0076723807, -0.009674717, -0.0002646696, 0.010251441, 0.013724749, 0.0076399804, 0.0009282674, 0.030637683, 0.012843462, 0.002658441, 0.0064703315, -0.007341898, 0.027112534, 0.0020865765, -0.017353578, 0.01001816, -0.0055728443, -0.00254504, 0.030482162, -0.016899973, 0.014683796, -0.00083916663, 0.014981879, -0.004305994, -0.01135953, 0.032944582, -0.0058741667, -0.010983687, -0.0042573935, -0.026101647, 0.010905926, -0.042820178, 0.0053687226, -0.007212297, -0.0019488754, -0.014787477, -0.0030002638, 0.005459443, 0.0019261952, -0.046500847, 0.026723731, -0.02529812, 0.028123423, -0.0155262025, -0.00754926, 0.022511698, 0.029756395, -0.018869909, 0.016887013, 0.005316882, -0.015565083, 0.017198056, 0.027112534, -0.004960479, -0.02023072, -0.029497193, 0.009577516, -0.012992503, -0.026801493, 0.0025887806, -0.024235392, 0.002924123, -0.015072599, -0.0033696266, 0.00883879, 0.006434691, -0.022058094, -0.01133361, 0.008339826, 0.018714389, -0.00084969675, 0.006868854, 0.015824284, -0.04548996, -0.01512444, 0.008592548, -0.16588931, 0.014606035, 0.03146713, -0.017379498, 0.016653731, 0.0011696493, 0.012150096, 0.004066232, -0.02156561, 0.0073548583, 0.008158385, 0.012078816, -0.04901511, -0.0010894586, 0.010225521, 0.019971518, -0.029756395, 0.014696756, 0.016614852, -0.0021805372, 0.028382625, 0.00043497345, -0.008547188, -0.015072599, -0.015435482, 0.0007221207, -0.0007225257, 0.036003165, 0.00630185, -0.018014543, -0.02760502, -0.0034279472, 0.018247824, 0.020671364, 0.0050512, -0.0059616473, -0.0066258525, 0.004617037, -0.005543684, 0.019803036, 0.015811326, 0.018714389, 0.005456203, 0.019155031, -0.0063342503, 0.02004928, 0.02014, -0.0021286968, -0.0125324195, -0.0120528955, 0.005427043, -0.042483218, 0.013478506, -0.0015325322, 0.020438083, 0.014230193, 0.014022831, 0.0004766888, 0.010815206, -0.0155262025, 0.0065027317, -0.010458803, 0.005556644, 0.019647516, -0.0010967486, -0.006107448, -0.0028107222, 0.011288249, 0.013284105, 0.007989903, -0.023107862, 0.0082750255, 0.004769318, -0.017301736, 0.01762574, -0.00377463, -0.018519986, -0.0038491504, 0.015915006, -0.0023668387, -0.015033719, 0.013284105, -0.003930151, 0.02515556, -0.007439099, 0.019258713, -0.0022210376, 0.010653204, -0.033437066, 0.0007978563, 0.013543307, -0.034836754, 0.018986551, -0.015979806, -0.025103718, 3.7361548e-05, 0.010977207, -0.008022304, -0.0015819425, 0.0023344385, -0.0136599485, -0.0036417888, -0.003956071, 0.0072511774, 0.024935238, -0.0006050748, -0.011152169, 0.0005568794, 0.023561466, -0.01884399, -0.016563011, 0.0030116038, 0.020425122, -0.013569227, 0.017327657, 0.021669291, 0.0076335003, -0.019854877, -0.017029574, -0.0073937387, 0.065941, -0.024520515, 0.010808726, -0.01892175, 0.014022831, -0.024196511, -0.081544966, -0.04162785, 0.019401273, 0.011502091, 0.011210489, 0.02396323, 0.005611724, 0.010737445, -0.026231248, 0.009681197, 0.008041744, -0.0077955015, 0.004163433, -0.004642957, 0.04162785, 0.0043189544, 0.011845534, -0.0032222054, -0.006169009, 0.012843462, -0.041057605, 0.007711261, -0.011657612, -0.027345816, -0.007678861, -0.012869382, -0.03434427, 0.028952869, 0.015552123, 0.010283842, -0.020788005, 0.0013672909, -0.0007164507, -0.009039671, -0.028253024, -0.006136609, -0.018507026, -0.03024888, 0.010724485, -0.03504412, 0.0003679454, 0.004140753, 0.003956071, -0.02659413, 0.003293486, 0.0032254455, -0.0029662435, 0.020749124, 0.002015296, -0.019219832, 0.008035264, 0.0059616473, -0.013005463, -0.030585842, 0.023380024, -0.0427165, 0.018338544, 0.014126512, -0.012506499, 0.0015049919, -0.012558339, 0.011761293, -0.007335418, 0.018325586, -0.009720077, -0.0042087934, -0.027060695, 0.007905663, 0.013828429, 0.011515051, -0.009305353, 0.014657876, -0.014826357, 0.0065999324, -0.030067438, -0.0008209415, -0.03268538, -0.0065999324, -0.016187169, -0.008281506, -0.019712316, -0.017146217, 0.010640244, -0.004380515, -0.004915119, 0.017560938, 0.0053395624, -0.016070528, -0.012778661, -0.0127916215, 0.018131183, 0.011832573, 0.025803564, 0.001023038, -0.0071280566, -0.005035, 0.012065856, 0.0013097804, 0.01388027, 0.013634028, -0.009804318, -0.013970991, -0.050025996, 0.035925403, -0.0039204312, -0.013096184, -0.009473835, -0.01621309, -0.006399051, -0.023017142, 0.003434427, 0.009052631, 0.012305617, -0.008657348, -0.024429793, -0.008547188, -0.015668765, -0.014359794, 0.011456731, -0.003648269, 0.0014450514, -0.00756222, 0.002656821, -0.022537619, -0.0006945805, 0.016887013, -0.007834382, 0.0143986745, -0.017172135, 0.00093960745, -0.008929511, -0.015228121, 0.015798366, -0.0252074, -0.0028155823, 0.015889086, -0.008229665, -0.031804092, -0.015785405, 0.011489131, 0.015785405, 0.0062370496, -0.003910711, -0.0356662, 0.006437931, -0.014644916, -0.02662005, 0.0056344047, -0.007108616, 0.015759485, 0.008430547, -0.021734092, 0.021513771, 0.022913462, -0.0038588706, 0.018584788, 0.003405267, -0.030767284, 0.006557812, -0.008411107, 0.006979015, -0.020775044, 0.0037713898, -0.002658441, 0.011469691, -0.020321442, 0.0075946203, 0.007704781, -0.011009607, -0.021073127, 0.00020483037, -0.018170064, 0.017133256, -0.00634397, 0.037195496, 0.013018423, 0.010685605, -0.014774517, -0.006139849, -0.0030504842, -0.025583243, 0.019349433, 0.009318314, -0.00081729644, -0.02165633, 0.029860076, 0.019608635, 0.028486306, -0.015565083, 0.0043448745, -0.0100959195, 0.020554723, -0.0029824437, 0.004529556, 0.0055987644, 0.02913431, -0.0019423953, 0.01648525, -0.0017998343, 0.02034736, -0.010361602, -0.001004408, -0.023056023, -0.007912142, -0.005958407, -0.0072252573, -0.019103192, 0.046552688, -0.010458803, -0.016575972, -0.00037786798, 0.027967902, 0.0016532231, -0.027319897, -0.023755869, 0.033799946, -0.027553178, -0.007989903, -0.0046073166, 0.0067262934, -0.018157104, 0.01137897, 0.016122367, 0.009745997, 0.025933165, -0.008378706, 0.006940135, 0.0023149983, 0.014139472, -0.032944582, 0.035795804, 0.01644637, 0.011165128, -0.021850733, -0.0025823005, -0.015189241, -0.0032627059, 0.0025515202, -0.021384168, 0.033592585, 0.0014734017, 0.050259277, 0.021876654, -0.035873566, 0.0070373355, -0.017547978, 0.012752741, -0.018636627, 0.0074909395, -0.016044607, 0.019699356, 0.009117432, 0.0017058735, -0.0027605018, -0.014476434, -0.035121877, 0.0033437065, -0.02519444, 0.014606035, -0.010769845, -0.03030072, 0.017936783, -0.011774253, 0.011508571, 0.014217232, -0.004950759, -0.033851787, 0.01905135, 0.038906228, -0.025946125, -0.008968391, -0.006081528, -0.0072835777, -0.019453114, -0.01902543, 0.013828429, -0.013024903, 0.0010359982, -0.038413744, -0.0042411936, 0.017016614, -0.02789014, 0.013711789, -0.024857476, -0.013439626, 0.025557322, -0.0036255887, -0.002768602, 0.0036417888, -0.00021384169}
)

func TestMilvusCollection(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		milvusClient, err := NewMilvusClient(addr, user, pwd)
		if err != nil {
			t.Fatal(err)
		}
		err = milvusClient.CreateCollection(context.Background(), fableDefaultCollection, 1)
		assert.Nil(t, err)
	})
	t.Run("drop", func(t *testing.T) {
		milvusClient, err := NewMilvusClient(addr, user, pwd)
		if err != nil {
			t.Fatal(err)
		}
		err = milvusClient.DropCollection(context.Background(), defaultCollection)
		assert.Nil(t, err)
	})
	t.Run("create index", func(t *testing.T) {
		milvusClient, err := NewMilvusClient(addr, user, pwd)
		if err != nil {
			t.Fatal(err)
		}
		idx, err := entity.NewIndexIvfFlat( // NewIndex func
			entity.L2, // metricType
			1024,      // ConstructParams
		)
		assert.Nil(t, err)
		err = milvusClient.CreateIndex(context.Background(), defaultCollection, defaultVectorColumn, idx, false)
		assert.Nil(t, err)
	})
	t.Run("insert", func(t *testing.T) {
		milvusClient, err := NewMilvusClient(addr, user, pwd)
		if err != nil {
			t.Fatal(err)
		}
		vec := make([][]float32, 0)
		vec = append(vec, inputVector)
		err = milvusClient.Insert(context.Background(), &model.InsertRequest{
			Collection: defaultCollection,
			ContentKey: filedContentKey,
			Content:    []string{"The dawn of Dark Mode"},
			Embeddings: vec,
			Metadata:   nil,
		})
		assert.Nil(t, err)
	})
	t.Run("ping", func(t *testing.T) {
		milvusClient, err := NewMilvusClient(addr, user, pwd)
		if err != nil {
			t.Fatal(err)
		}
		_, err = milvusClient.Search(context.Background(), &model.SearchRequest{
			Collection: "",
			Vector:     nil,
		})
		assert.Nil(t, err)
	})
	t.Run("search", func(t *testing.T) {
		milvusClient, err := NewMilvusClient(addr, user, pwd)
		if err != nil {
			t.Fatal(err)
		}
		inputVector = []float32{0.014838364, -0.017620698, 0.039551493, 0.015700748, -0.0011719975, -0.013021858, -0.010251275, 0.012758525, 0.019684473, -0.03831288, 0.001341002, 0.064250916, -0.002414946, 0.013822225, -0.0031436041, 0.019625148, -0.01564217, -0.0052712923, -0.015714055, -0.030302437, -0.023091216, -0.009670714, 0.015752068, -0.0038696309, -0.0043562334, -0.008719393, 0.0027820505, 0.02187728, 0.020052845, -0.0077083246, -0.0087030055, -0.008234454, 0.0040822765, -0.015165803, 0.00357362, 0.015885333, 0.023709305, -0.014285445, -0.018217837, -0.0178895, -0.036025092, 0.00083955115, 0.014567554, 0.005873727, 0.013601852, 0.013569581, -0.13623272, -0.01620436, -0.016131273, -0.036731627, 0.021327132, 0.00086528336, 0.027448554, 0.06727927, 0.014286099, 0.013202444, -0.02393746, 0.0235503, 0.019170485, 0.021919766, 0.03503096, -0.019138088, -0.006739168, 0.027075253, -0.033817433, 0.036604814, -0.010031421, -0.016051564, 0.0076677925, -0.009989969, 0.009421281, 0.0050558955, -0.013369248, -0.013520031, -0.021161506, 0.0032631315, -0.002559837, -0.008193062, 0.03864989, -0.0041056178, 0.011382851, 0.0043259826, 0.044303134, 0.0062390473, -0.007576711, 0.0070692142, 0.011444597, 0.002004812, 0.011841655, 0.03266881, -0.0440562, -0.019674443, 0.049414042, 0.03268112, 0.006078411, -0.017142553, 0.020702785, 0.007652628, 0.020619985, 0.009974223, 0.007835002, 0.04429414, -0.0064316993, -0.052523065, -0.029133398, -0.0146194985, -0.014619836, 0.009158691, 0.008382959, -0.10954855, -0.019274492, 0.030559044, -0.02578594, 0.018183038, -0.03630349, 0.028039437, -0.0031390805, 0.0137243755, 0.030133571, 0.0005093897, 0.008661845, 0.05171257, 0.033242453, 0.013949195, 0.020825248, 0.030746315, 0.014621836, -0.0046388065, -0.030364914, 0.01363981, 0.013196568, 0.030295568, -0.0174938, -0.020367127, 0.027555155, -0.018687103, 0.051717162, 0.023337685, 0.0068090595, 0.031118283, -0.04470329, -0.03815111, -0.16577286, -0.009473712, 0.0012697432, -0.035571016, 0.015115573, -0.019996492, -0.03263, 0.011677124, 0.01732706, -0.00072131, -0.032654494, 0.019310793, 0.005825045, 0.0024686672, 0.016264752, -0.0080637615, 0.0099125365, 0.02263525, -0.023083536, 0.012654782, -0.0019935432, 0.0013407182, 0.033267837, -0.01659477, -0.0017987847, -0.0115283, 0.037528917, 0.00034592906, -0.026666755, -0.017653765, -0.011911211, -0.044888854, 0.010300219, 0.017710255, 0.06358, 0.018672366, -0.021409824, -0.023174927, -0.024199752, -0.0200502, -0.005743808, -0.005195007, -0.011057331, -0.01153751, 0.006358018, 0.020102525, 0.0036482087, 0.012744272, -0.010941555, 0.03139748, 0.040872224, -0.023371434, 0.04934634, -0.032586604, -0.029276755, 0.040274605, 0.012786386, 0.029377582, 0.009591064, 0.00343692, 0.015113719, -0.0059439307, 0.012573738, 0.15877238, -0.011829165, 0.0047420408, -0.0033891946, 0.0073359436, 0.01082526, -0.011793815, -0.016057836, 0.020213237, 0.028924668, -0.0030009332, -0.03922341, -0.008591165, 0.0137856705, -0.008285535, 0.023076719, -0.0039982805, -0.0146758715, -0.022576937, -0.015657607, 0.024632135, 0.03174942, 0.010056781, 0.0027731599, -0.06140183, -0.01037124, 0.0016801344, -0.038836733, 0.028363021, 0.02272064, -0.0031611838, 0.0038485106, -0.021977214, 0.02571117, -0.044925276, 0.0065538604, -0.038620915, 0.032227237, 0.04337141, -0.004964503, 0.03319142, -0.0025206713, 0.03564532, -0.03317853, -0.016576312, -0.010151547, -0.019243194, 0.01859985, -0.011652391, 0.011788655, -0.008821268, -0.013966298, -0.029581772, 0.0011792006, 0.022502907, 0.023570362, 0.011041806, 0.0054137483, 0.01136796, 0.025268491, -0.025972445, -0.033805266, -0.04554501, 0.04532494, 0.019158257, -0.005576231, 0.031196062, -0.015332367, -0.22690581, -0.018452574, 0.0028970218, 0.028916495, 0.036128055, 0.02234652, 0.035011537, 0.011490907, 0.0154022435, 0.034814894, 0.016039088, -0.033290755, 0.009284299, 0.0345211, -0.041654, 0.013161275, 0.00090003444, -0.0066361027, 0.014908384, -0.0007631966, 0.012913064, 0.0036558479, 0.0005156213, -0.0070527974, -0.004504945, 0.0029339797, 0.010207012, -0.007562483, 0.029348409, -0.048445627, -0.00320237, -0.014569097, -0.04028096, -0.016391898, 0.02919993, -0.47404745, -0.013345222, -0.022310713, -0.026250625, -0.01452949, -0.011486279, -0.0032556157, 0.019542836, 0.0001164813, -0.006093245, -0.0014315548, 0.019292241, 0.002716208, 0.009988825, -0.028337467, -0.001480295, -0.010886474, 0.01840926, 0.006160934, -0.00936346, -0.03437453, -0.024849296, -0.0020006124, -0.013547301, -0.017221741, 0.004963671, -0.0058589587, 0.010824432, 0.002591003, 0.010618487, 0.004222109, -0.003241194, 0.016680053, 0.013110307, 0.014485241, 0.0044939877, 0.026941715, 0.00865972, 0.03758676, 0.0072575766, -0.014026135, 0.059755262, 0.0049331742, 0.024177276, 0.040791262, 0.034922283, 0.0017641912, 0.0076271654, -0.035317704, 0.02381471, -0.01385649, -0.007852505, 0.041393764, -0.0050646984, 0.011507972, 0.012655169, -0.0076508056, -0.019883038, 0.017577171, -0.0002979984, 0.0031304385, -0.008362382, 0.020301871, -0.019696226, 0.040463403, -0.0005293077, 0.009099654, -0.019911662, 0.048317768, 0.035024546, -0.017664244, 0.0019191966, 0.002901581, -0.08070902, 0.00021889158, 0.021420196, -0.034751747, -0.019428387, 0.03938823, -0.024003804, -0.0046148803, -0.016937349, 0.032183547, 0.00402007, 0.028713647, -0.017370513, 0.019817289, -0.02342447, 0.0068543674, 0.00024666477, -0.00006084482, -0.039199762, -0.018024683, -0.03164163, 0.051410977, 0.048208483, -0.0119695775, -0.026077947, 0.029647399, 0.0028355687, 0.034863584, -0.016725319, -0.0038504696, -0.0031427776, -0.028767373, -0.01213128, 0.002139904, 0.01762601, 0.006600361, -0.053048335, -0.020808456, 0.011661664, -0.036584266, -0.0020801786, 0.06361078, -0.0135586, -0.023307404, 0.0067742546, 0.0036768273, 0.008805854, -0.014147713, 0.024371535, 0.011926635, -0.022567747, -0.0012365249, 0.0087882625, -0.028221928, 0.0022280738, 0.0034959784, 0.02975402, 0.0108348485, 0.011475353, 0.032121204, -0.050673988, -0.0026767442, -0.04551726, -0.029637087, 0.011718687, -0.006415919, -0.003757314, -0.020736735, 0.016643772, 0.025524274, 0.03153301, 0.04350964, -0.0036638333, 0.022227088, -0.022370365, -0.03546318, 0.016464086, 0.02993969, -0.028887564, 0.061588686, 0.0039404524, -0.021679886, 0.0015329276, 0.03013087, -0.022506868, -0.018775234, -0.014957843, 0.0119680045, -0.017912662, -0.08536765, 0.011492664, 0.025277797, 0.028616965, -0.014308099, -0.056252003, 0.005479866, -0.0048354478, -0.029662022, -0.0073950156, -0.021931097, -0.006393613, 0.010457335, 0.016776003, -0.022647167, -0.016725188, 0.007636527, -0.037838165, 0.017805288, 0.016914202, -0.028953755, 0.0027893241, 0.020019934, -0.027045246, -0.0063255937, -0.04483596, 0.00076982466, 0.018942785, 0.031327553, 0.008128298, 0.020161482, -0.021124803, -0.024175653, -0.013439824, -0.00020664975, 0.047930967, 0.002983051, -0.0064846002, 0.0014040401, 0.010465845, -0.036296282, 0.004096728, 0.045669038, -0.034743734, 0.054350525, 0.0024032863, -0.02436844, 0.0021644612, -0.0015580772, -0.030310983, -0.02116161, 0.004089734, 0.012532454, -0.011866499, -0.008149107, -0.015987147, -0.014838581, -0.014042834, 0.03950637, -0.0081894165, -0.011330618, -0.0044439165, -0.017552774, 0.0055255364, 0.047746133, -0.0180145, -0.006038069, 0.026872234, -0.0022477505, 0.013040733, 0.023238622, -0.0010938424, 0.003337597, -0.013552453, -0.014337369, -0.0017943259, -0.00930265, 0.034664348, 0.019636748, 0.018359799, -0.012646118, -0.03573063, -0.023170844, 0.022474239, -0.009263427, -0.008184437, -0.02938159, 0.014352807, 0.010041592, -0.023655923, -0.012420379, -0.039050147, 0.044157572, -0.003995087, 0.00902031, 0.0054842355, -0.0033872214, -0.011662388, -0.037438367, 0.013023613, -0.014543291, -0.0072621885, -0.005672118, 0.03369685, -0.0027222203, -0.014462151, -0.008329522, 0.009433829, 0.0136858905, -0.024328247, 0.006465071, 0.020336792, -0.008821166, -0.0065075024, -0.04063474, 0.008280466, -0.0015778339, -0.0023045647, 0.11269022, 0.015746169, -0.024842395, 0.00055270403, 0.009698641, -0.000575511, 0.024415608, -0.019155748, 0.021760643, 0.0069623054, 0.032219592, -0.0020632427, 0.020025678, -0.0007373515, 0.017476436, 0.053086706, 0.022838505, -0.0139748575, 0.024216855, -0.038033184, 0.026937954, 0.043613244, 0.03662654, -0.027365856, 0.008198622, -0.012873694, 0.010881722, -0.0040878137, -0.038802564, -0.0040000193, 0.067471705, 0.032404773, -0.018721236, 0.060300224, -0.011354774, 0.021440435, 0.045396443, -0.016600858, -0.010134071, -0.002838615, 0.027490908, 0.028651439, 0.025606519, -0.008494687, -0.017257677, 0.013411966, 0.0051182695, 0.013201371, 0.010941291, -0.012172995, 0.042408634, 0.011018278, -0.014006068, 0.048964094, -0.025562942, -0.022584576, 0.049159188, -0.037635382, -0.032043677, 0.013482845, 0.06878712, 0.038469125, 0.0030257653, 0.0033526968, 0.0066778027, 0.031604078, 0.019530509, 0.030301707, -0.007821905, -0.010145646, 0.011003401, -0.028458841, 0.0051495847, -0.026235644, -0.0008209383, -0.029649906, 0.009109153, -0.03184981, -0.0041207597, -0.20402302, 0.0051892903, 0.010797182, 0.02067622, 0.018927343, 0.009714799, -0.006196521, -0.011524141, 0.0023700763, 0.012265317, -0.036812942, -0.009687532, 0.014553207, 0.020836778, -0.022368807, -0.04392204, -0.017737128, 0.021430675, 0.013592231, 0.012731836, 0.0038867812, -0.03333643, -0.03177168, 0.018369831, -0.0044153817, 0.023774516, 0.0033385134, 0.009527718, 0.0032199968, -0.008123785, -0.010461502, -0.020104649, -0.010899571, 0.0036919126, -0.01505985, -0.035547256, 0.029195618, 0.00836193, -0.0031257174, -0.012548003, 0.008276347, -0.0062073655, 0.008377875, 0.0044152043, 0.009370235, -0.010262294, 0.026213396, -0.047776476, -0.0076019284, -0.046072166, 0.034080144, 0.04708474, 0.016887741, -0.009885441, 0.030334515, -0.017878516, 0.018351035, 0.024792878, -0.015724272, 0.008981331, -0.018996129, -0.01979231, 0.020773347, 0.0126279155, 0.0044938116, -0.01187637, -0.0005444625, 0.006843281, 0.023123201, 0.020533776, 0.024486467, -0.015983341, -0.026969647, 0.0023893982, 0.0315646, -0.02402019, 0.0054192157, 0.027887626, 0.0048745675, 0.024269693, -0.0074570943, -0.028718008, -0.0071378606, -0.0058130594, -0.0017354892, -0.42743066, -0.01868335, -0.02150815, -0.014938819, -0.018707903, -0.0062781246, -0.016588869, -0.016222054, 0.015956137, -0.0066683292, -0.0014411178, -0.021973483, 0.02273645, -0.021794233, 0.016480051, 0.006484712}
		res, err := milvusClient.Search(context.Background(), &model.SearchRequest{
			Collection: "",
			Vector:     nil,
		})
		assert.Nil(t, err)
		for _, i2 := range res {
			fmt.Printf("search result val:%v score:%v\n", i2.Payload, i2.Score)
		}
	})
}
