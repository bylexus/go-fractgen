package lib

import (
	"math"
	"math/big"
)

const SYS_PRECISION = 128

var BIG_2 = big.NewFloat(2.0).SetPrec(SYS_PRECISION)
var BIG_3 = big.NewFloat(3.0).SetPrec(SYS_PRECISION)
var BIG_4 = big.NewFloat(4.0).SetPrec(SYS_PRECISION)
var BIG_6 = big.NewFloat(6.0).SetPrec(SYS_PRECISION)

var BIG_MAX_ABS_SQUARE_AMOUNT = big.NewFloat(256).SetPrec(SYS_PRECISION)
var BIG_LOG_MAX_ABS_SQUARE_AMOUNT *big.Float
var BIG_LOG_2 *big.Float

func init() {
	BIG_LOG_MAX_ABS_SQUARE_AMOUNT = big.NewFloat(math.Log(256.0)).SetPrec(SYS_PRECISION)
	BIG_LOG_2 = big.NewFloat(math.Log(2.0)).SetPrec(SYS_PRECISION)
}
