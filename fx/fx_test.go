package fx

import (
	"strings"
	"testing"

	"github.com/shopspring/decimal"
)

const sample = `<?xml version="1.0" encoding="UTF-8"?>
<gesmes:Envelope xmlns:gesmes="http://www.gesmes.org/xml/2002-08-01" xmlns="http://www.ecb.int/vocabulary/2002-08-01/eurofxref">
	<Cube>
		<Cube time='2026-10-06'>
			<Cube currency='USD' rate='1.1269'/>
			<Cube currency='GBP' rate='0.84880'/>
		</Cube>
	</Cube>
</gesmes:Envelope>`

func TestToEUR(t *testing.T) {
	rates, err := parseECB(strings.NewReader(sample))
	if err != nil {
		t.Fatal(err)
	}
	got, err := rates.ToEUR(decimal.RequireFromString("84.88"), "GBP")
	if err != nil || !got.Equal(decimal.NewFromInt(100)) {
		t.Errorf("84.88 GBP = %v, %v; want 100", got, err)
	}
	if got, _ := rates.ToEUR(decimal.NewFromInt(5), ""); !got.Equal(decimal.NewFromInt(5)) {
		t.Errorf("empty currency should be EUR, got %v", got)
	}
	if _, err := rates.ToEUR(decimal.NewFromInt(5), "XYZ"); err == nil {
		t.Error("expected error for unknown currency")
	}
}

func TestParseFallback(t *testing.T) {
	rates, err := parseFallback(strings.NewReader(`{"result":"success","rates":{"EUR":1,"RUB":96.5302,"KZT":509.35}}`))
	if err != nil {
		t.Fatal(err)
	}
	got, err := rates.ToEUR(decimal.RequireFromString("965.302"), "RUB")
	if err != nil || !got.Equal(decimal.NewFromInt(10)) {
		t.Errorf("965.302 RUB = %v, %v; want 10", got, err)
	}
	if _, err := parseFallback(strings.NewReader(`{"result":"error"}`)); err == nil {
		t.Error("expected error for failed result")
	}
}
