package types

import (
	"fmt"

	sdkmath "cosmossdk.io/math"
	cosmos "github.com/cosmos/cosmos-sdk/types"
	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/ethereum/go-ethereum/params"
)

type (
	LegacyParams = paramtypes.ParamSet
	// Subspace defines an interface that implements the legacy Cosmos SDK x/params Subspace type.
	// NOTE: This is used solely for migration of the Cosmos SDK x/params managed parameters.
	Subspace interface {
		GetParamSetIfExists(ctx cosmos.Context, ps LegacyParams)
	}
)

var _ paramtypes.ParamSet = (*Params)(nil)

var (
	// DefaultMinGasMultiplier is 0.5 or 50%
	DefaultMinGasMultiplier = sdkmath.LegacyNewDecWithPrec(50, 2)
	// DefaultMinGasPrice is 20Gwei (i.e disabled)
	DefaultMinGasPrice = sdkmath.LegacyNewDec(20000000000)
	// DefaultEnableHeight is 0 (i.e disabled)
	DefaultEnableHeight = int64(0)
	// DefaultNoBaseFee is false
	DefaultNoBaseFee = false
)

// Parameter keys
var (
	ParamStoreKeyNoBaseFee                = []byte("NoBaseFee")
	ParamStoreKeyBaseFeeChangeDenominator = []byte("BaseFeeChangeDenominator")
	ParamStoreKeyElasticityMultiplier     = []byte("ElasticityMultiplier")
	ParamStoreKeyBaseFee                  = []byte("BaseFee")
	ParamStoreKeyEnableHeight             = []byte("EnableHeight")
	ParamStoreKeyMinGasPrice              = []byte("MinGasPrice")
	ParamStoreKeyMinGasMultiplier         = []byte("MinGasMultiplier")
)

// ParamKeyTable the param key table for launch module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// NewParams creates a new Params instance
func NewParams(
	noBaseFee bool,
	baseFeeChangeDenom,
	elasticityMultiplier uint32,
	baseFee uint64,
	enableHeight int64,
	minGasPrice sdkmath.LegacyDec,
	minGasPriceMultiplier sdkmath.LegacyDec,
) Params {
	return Params{
		NoBaseFee:                noBaseFee,
		BaseFeeChangeDenominator: baseFeeChangeDenom,
		ElasticityMultiplier:     elasticityMultiplier,
		BaseFee:                  sdkmath.NewIntFromUint64(baseFee),
		EnableHeight:             enableHeight,
		MinGasPrice:              minGasPrice,
		MinGasMultiplier:         minGasPriceMultiplier,
	}
}

// DefaultParams returns a default set of parameters
func DefaultParams() Params {
	return Params{
		NoBaseFee:                DefaultNoBaseFee,
		BaseFeeChangeDenominator: params.DefaultBaseFeeChangeDenominator,
		ElasticityMultiplier:     params.DefaultElasticityMultiplier,
		BaseFee:                  sdkmath.NewIntFromUint64(params.InitialBaseFee),
		EnableHeight:             DefaultEnableHeight,
		MinGasPrice:              DefaultMinGasPrice,
		MinGasMultiplier:         DefaultMinGasMultiplier,
	}
}

// ParamSetPairs get the params.ParamSet
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(ParamStoreKeyNoBaseFee, &p.NoBaseFee, validateBool),
		paramtypes.NewParamSetPair(ParamStoreKeyBaseFeeChangeDenominator, &p.BaseFeeChangeDenominator, validateBaseFeeChangeDenominator),
		paramtypes.NewParamSetPair(ParamStoreKeyElasticityMultiplier, &p.ElasticityMultiplier, validateElasticityMultiplier),
		paramtypes.NewParamSetPair(ParamStoreKeyBaseFee, &p.BaseFee, validateBaseFee),
		paramtypes.NewParamSetPair(ParamStoreKeyEnableHeight, &p.EnableHeight, validateEnableHeight),
		paramtypes.NewParamSetPair(ParamStoreKeyMinGasPrice, &p.MinGasPrice, validateMinGasPrice),
		paramtypes.NewParamSetPair(ParamStoreKeyMinGasMultiplier, &p.MinGasMultiplier, validateMinGasPrice),
	}
}

// Validate validates the set of params
func (p Params) Validate() error {
	if p.BaseFeeChangeDenominator == 0 {
		return fmt.Errorf("base fee change denominator cannot be 0")
	}

	if p.BaseFee.IsNegative() {
		return fmt.Errorf("initial base fee cannot be negative: %s", p.BaseFee)
	}

	if p.EnableHeight < 0 {
		return fmt.Errorf("enable height cannot be negative: %d", p.EnableHeight)
	}

	if err := validateMinGasMultiplier(p.MinGasMultiplier); err != nil {
		return err
	}

	return validateMinGasPrice(p.MinGasPrice)
}

func validateBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func (p *Params) IsBaseFeeEnabled(height int64) bool {
	return !p.NoBaseFee && height >= p.EnableHeight
}

func validateMinGasPrice(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("value cannot be negative: %s", i)
	}

	return nil
}

func validateBaseFeeChangeDenominator(i interface{}) error {
	value, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value == 0 {
		return fmt.Errorf("base fee change denominator cannot be 0")
	}

	return nil
}

func validateElasticityMultiplier(i interface{}) error {
	_, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}
	return nil
}

func validateBaseFee(i interface{}) error {
	value, ok := i.(sdkmath.Int)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value.IsNegative() {
		return fmt.Errorf("base fee cannot be negative")
	}

	return nil
}

func validateEnableHeight(i interface{}) error {
	value, ok := i.(int64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if value < 0 {
		return fmt.Errorf("enable height cannot be negative: %d", value)
	}

	return nil
}

func validateMinGasMultiplier(i interface{}) error {
	v, ok := i.(sdkmath.LegacyDec)

	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v.IsNil() {
		return fmt.Errorf("invalid parameter: nil")
	}

	if v.IsNegative() {
		return fmt.Errorf("value cannot be negative: %s", v)
	}

	if v.GT(sdkmath.LegacyOneDec()) {
		return fmt.Errorf("value cannot be greater than 1: %s", v)
	}
	return nil
}
