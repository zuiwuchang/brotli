package brotli

import "math"

/* Copyright 2013 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Literal cost model to allow backward reference replacement to be efficient.
 */
/* Copyright 2013 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Literal cost model to allow backward reference replacement to be efficient.
 */

/* Estimates how many bits the literals in the interval [pos, pos + len) in the
   ring-buffer (data, mask) will take entropy coded and writes these estimates
   to the cost[0..len) array. */
/* Copyright 2013 Google Inc. All Rights Reserved.

   Distributed under MIT license.
   See file LICENSE for detail or copy at https://opensource.org/licenses/MIT
*/

/* Utilities for fast computation of logarithms. */
func log2FloorNonZero(n uint) uint32 {
	/* TODO: generalize and move to platform.h */
	var result uint32 = 0
	for {
		n >>= 1
		if n == 0 {
			break
		}
		result++
	}
	return result
}

/* A lookup table for small values of log2(int) to be used in entropy
   computation.

   ", ".join(["%.16ff" % x for x in [0.0]+[log2(x) for x in range(1, 256)]]) */
var kLog2Table = []float32{
	0.0000000000000000,
	0.0000000000000000,
	1.0000000000000000,
	1.5849625007211563,
	2.0000000000000000,
	2.3219280948873622,
	2.5849625007211561,
	2.8073549220576042,
	3.0000000000000000,
	3.1699250014423126,
	3.3219280948873626,
	3.4594316186372978,
	3.5849625007211565,
	3.7004397181410922,
	3.8073549220576037,
	3.9068905956085187,
	4.0000000000000000,
	4.0874628412503400,
	4.1699250014423122,
	4.2479275134435852,
	4.3219280948873626,
	4.3923174227787607,
	4.4594316186372973,
	4.5235619560570131,
	4.5849625007211570,
	4.6438561897747244,
	4.7004397181410926,
	4.7548875021634691,
	4.8073549220576037,
	4.8579809951275728,
	4.9068905956085187,
	4.9541963103868758,
	5.0000000000000000,
	5.0443941193584534,
	5.0874628412503400,
	5.1292830169449664,
	5.1699250014423122,
	5.2094533656289501,
	5.2479275134435852,
	5.2854022188622487,
	5.3219280948873626,
	5.3575520046180838,
	5.3923174227787607,
	5.4262647547020979,
	5.4594316186372973,
	5.4918530963296748,
	5.5235619560570131,
	5.5545888516776376,
	5.5849625007211570,
	5.6147098441152083,
	5.6438561897747244,
	5.6724253419714961,
	5.7004397181410926,
	5.7279204545631996,
	5.7548875021634691,
	5.7813597135246599,
	5.8073549220576046,
	5.8328900141647422,
	5.8579809951275719,
	5.8826430493618416,
	5.9068905956085187,
	5.9307373375628867,
	5.9541963103868758,
	5.9772799234999168,
	6.0000000000000000,
	6.0223678130284544,
	6.0443941193584534,
	6.0660891904577721,
	6.0874628412503400,
	6.1085244567781700,
	6.1292830169449672,
	6.1497471195046822,
	6.1699250014423122,
	6.1898245588800176,
	6.2094533656289510,
	6.2288186904958804,
	6.2479275134435861,
	6.2667865406949019,
	6.2854022188622487,
	6.3037807481771031,
	6.3219280948873617,
	6.3398500028846252,
	6.3575520046180847,
	6.3750394313469254,
	6.3923174227787598,
	6.4093909361377026,
	6.4262647547020979,
	6.4429434958487288,
	6.4594316186372982,
	6.4757334309663976,
	6.4918530963296748,
	6.5077946401986964,
	6.5235619560570131,
	6.5391588111080319,
	6.5545888516776376,
	6.5698556083309478,
	6.5849625007211561,
	6.5999128421871278,
	6.6147098441152092,
	6.6293566200796095,
	6.6438561897747253,
	6.6582114827517955,
	6.6724253419714952,
	6.6865005271832185,
	6.7004397181410917,
	6.7142455176661224,
	6.7279204545631988,
	6.7414669864011465,
	6.7548875021634691,
	6.7681843247769260,
	6.7813597135246599,
	6.7944158663501062,
	6.8073549220576037,
	6.8201789624151887,
	6.8328900141647422,
	6.8454900509443757,
	6.8579809951275719,
	6.8703647195834048,
	6.8826430493618416,
	6.8948177633079437,
	6.9068905956085187,
	6.9188632372745955,
	6.9307373375628867,
	6.9425145053392399,
	6.9541963103868758,
	6.9657842846620879,
	6.9772799234999168,
	6.9886846867721664,
	7.0000000000000000,
	7.0112272554232540,
	7.0223678130284544,
	7.0334230015374501,
	7.0443941193584534,
	7.0552824355011898,
	7.0660891904577721,
	7.0768155970508317,
	7.0874628412503400,
	7.0980320829605272,
	7.1085244567781700,
	7.1189410727235076,
	7.1292830169449664,
	7.1395513523987937,
	7.1497471195046822,
	7.1598713367783891,
	7.1699250014423130,
	7.1799090900149345,
	7.1898245588800176,
	7.1996723448363644,
	7.2094533656289492,
	7.2191685204621621,
	7.2288186904958804,
	7.2384047393250794,
	7.2479275134435861,
	7.2573878426926521,
	7.2667865406949019,
	7.2761244052742384,
	7.2854022188622487,
	7.2946207488916270,
	7.3037807481771031,
	7.3128829552843557,
	7.3219280948873617,
	7.3309168781146177,
	7.3398500028846243,
	7.3487281542310781,
	7.3575520046180847,
	7.3663222142458151,
	7.3750394313469254,
	7.3837042924740528,
	7.3923174227787607,
	7.4008794362821844,
	7.4093909361377026,
	7.4178525148858991,
	7.4262647547020979,
	7.4346282276367255,
	7.4429434958487288,
	7.4512111118323299,
	7.4594316186372973,
	7.4676055500829976,
	7.4757334309663976,
	7.4838157772642564,
	7.4918530963296748,
	7.4998458870832057,
	7.5077946401986964,
	7.5156998382840436,
	7.5235619560570131,
	7.5313814605163119,
	7.5391588111080319,
	7.5468944598876373,
	7.5545888516776376,
	7.5622424242210728,
	7.5698556083309478,
	7.5774288280357487,
	7.5849625007211561,
	7.5924570372680806,
	7.5999128421871278,
	7.6073303137496113,
	7.6147098441152075,
	7.6220518194563764,
	7.6293566200796095,
	7.6366246205436488,
	7.6438561897747244,
	7.6510516911789290,
	7.6582114827517955,
	7.6653359171851765,
	7.6724253419714952,
	7.6794800995054464,
	7.6865005271832185,
	7.6934869574993252,
	7.7004397181410926,
	7.7073591320808825,
	7.7142455176661224,
	7.7210991887071856,
	7.7279204545631996,
	7.7347096202258392,
	7.7414669864011465,
	7.7481928495894596,
	7.7548875021634691,
	7.7615512324444795,
	7.7681843247769260,
	7.7747870596011737,
	7.7813597135246608,
	7.7879025593914317,
	7.7944158663501062,
	7.8008998999203047,
	7.8073549220576037,
	7.8137811912170374,
	7.8201789624151887,
	7.8265484872909159,
	7.8328900141647422,
	7.8392037880969445,
	7.8454900509443757,
	7.8517490414160571,
	7.8579809951275719,
	7.8641861446542798,
	7.8703647195834048,
	7.8765169465650002,
	7.8826430493618425,
	7.8887432488982601,
	7.8948177633079446,
	7.9008668079807496,
	7.9068905956085187,
	7.9128893362299619,
	7.9188632372745955,
	7.9248125036057813,
	7.9307373375628867,
	7.9366379390025719,
	7.9425145053392399,
	7.9483672315846778,
	7.9541963103868758,
	7.9600019320680806,
	7.9657842846620870,
	7.9715435539507720,
	7.9772799234999168,
	7.9829935746943104,
	7.9886846867721664,
	7.9943534368588578,
}

/* Faster logarithm for small integers, with the property of log2(0) == 0. */
func fastLog2(v uint) float64 {
	if v < uint(len(kLog2Table)) {
		return float64(kLog2Table[v])
	}

	return math.Log2(float64(v))
}
