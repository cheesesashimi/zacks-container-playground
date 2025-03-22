package genflag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBoolOptionFuncs(t *testing.T) {
	bf := boolFlag{value: false}

	assert.NoError(t, applyOptionFuncs(&bf, Single, Quoted, EqualSeparator, Explicit, Title))

	expected := boolFlag{
		flagOpts: flagOpts{
			single:    true,
			quoted:    true,
			separator: "=",
		},
		explicit: true,
	}

	assert.Equal(t, expected.flagOpts, bf.flagOpts)
	assert.Equal(t, expected.explicit, bf.explicit)
	assert.NotNil(t, bf.transform)
}

func TestStringOptionFuncs(t *testing.T) {
	sf := stringFlag{}

	expected := stringFlag{
		flagOpts: flagOpts{
			single:    true,
			quoted:    true,
			separator: " ",
		},
	}

	assert.NoError(t, applyOptionFuncs(&sf, Single, Quoted, Separator(" ")))

	assert.Equal(t, expected, sf)

	assert.Error(t, applyOptionFuncs(&sf, Single, Quoted, Separator("="), Explicit))
}
