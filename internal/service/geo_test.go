package service

import (
	"context"
	"fmt"
	"os"
	"testing"

	mocks "vio_coding_challenge/mocks/repo"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TODO It would be good to add benchmarks on a real CSV file and a real database through testing.B

func Test_GeoService_ParseCSV_Ok(t *testing.T) {
	filePath := os.TempDir() + "/Test_GeoService_ParseCSV_Ok.txt"

	defer func() {
		if dErr := os.Remove(filePath); dErr != nil {
			t.Errorf("Failed to remove file %q: %v\n", filePath, dErr)
		}
	}()

	cases := []struct {
		name    string
		data    string
		stats   ParsingStats
		callsFn func(m *mocks.MockGeoRepoI)
	}{
		{
			name: "happy path",
			data: `ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,68.31023296602508,-37.62435199624531,7301823115
70.95.73.73,TL,Saudi Arabia,Gradymouth,-49.16675918861615,-86.05920084416894,2559997162
,PY,Falkland Islands (Malvinas),,75.41685191518815,-144.6943217219469,0
125.159.20.54,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276`,
			stats: ParsingStats{
				RecAcceptedCnt:  4,
				RecDiscardedCnt: 1,
			},
			callsFn: func(m *mocks.MockGeoRepoI) {
				m.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(int64(2), nil)
				m.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(int64(2), nil)
			},
		},
		{
			name: "all data are accepted",
			data: `ip_address,country_code,country,city,latitude,longitude,mystery_value
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115
125.159.20.54,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276`,
			stats: ParsingStats{
				RecAcceptedCnt:  3,
				RecDiscardedCnt: 0,
			},
			callsFn: func(m *mocks.MockGeoRepoI) {
				m.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(int64(2), nil)
				m.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
		},
		{
			name: "empty file",
			data: `ip_address,country_code,country,city,latitude,longitude,mystery_value`,
			stats: ParsingStats{
				RecAcceptedCnt:  0,
				RecDiscardedCnt: 0,
			},
			callsFn: func(m *mocks.MockGeoRepoI) {},
		},
		{
			name: "check not null constraints",
			data: `ip_address,country_code,country,city,latitude,longitude,mystery_value
,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
125.159.20.54,,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
125.159.20.54,LI,,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
125.159.20.54,LI,Guyana,,-78.2274228596799,-163.26218895343357,1337885276
125.159.20.54,LI,Guyana,Port Karson,,-163.26218895343357,1337885276
125.159.20.54,LI,Guyana,Port Karson,-78.2274228596799,,1337885276
125.159.20.54,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,
,,,,,,`,
			stats: ParsingStats{
				RecAcceptedCnt:  0,
				RecDiscardedCnt: 8,
			},
			callsFn: func(m *mocks.MockGeoRepoI) {},
		},
		{
			name: "check ip address constraint",
			data: `ip_address,country_code,country,city,latitude,longitude,mystery_value
,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
wrong type,,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
1251592054,LI,,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
125.159.20,LI,Guyana,,-78.2274228596799,-163.26218895343357,1337885276
125 159 20 54,LI,Guyana,Port Karson,,-163.26218895343357,1337885276
125.159.20.54,LI,Guyana,Port Karson,-78.2274228596799,78.2274228596799,1337885276
0.0.0.0,LI,Guyana,Port Karson,-78.2274228596799,78.2274228596799,1337885276
localhost,LI,Guyana,Port Karson,-78.2274228596799,78.2274228596799,1337885276`,
			stats: ParsingStats{
				RecAcceptedCnt:  2,
				RecDiscardedCnt: 6,
			},
			callsFn: func(m *mocks.MockGeoRepoI) {
				m.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(int64(2), nil)
			},
		},
		{
			name: "check country code constraint",
			data: `ip_address,country_code,country,city,latitude,longitude,mystery_value
123.123.123.123,LI,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
123.123.123.123,123,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
123.123.123.123,li,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
123.123.123.123,qweqwewqwe,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276
123.123.123.123,,Guyana,Port Karson,-78.2274228596799,-163.26218895343357,1337885276`,
			stats: ParsingStats{
				RecAcceptedCnt:  1,
				RecDiscardedCnt: 4,
			},
			callsFn: func(m *mocks.MockGeoRepoI) {
				m.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
		},
		{
			name: "check size constraint",
			data: "ip_address,country_code,country,city,latitude,longitude,mystery_value\n" +
				"123.123.123.123,LI,TESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTT" +
				"ESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTT" +
				"ESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTT" +
				"ESTTESTTESTTESTTESTTESTTESTTESTTESTTEST,Port Karson,-78.2274228596799,-163.26218895343357,1337885" +
				"276\n" +
				"123.123.123.123,LI,Guyana,TESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTE" +
				"STTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTE" +
				"STTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTE" +
				"STTESTTESTTESTTESTTESTTESTTESTTESTTESTTESTTEST,-78.2274228596799,-163.26218895343357,1337885276",
			stats: ParsingStats{
				RecAcceptedCnt:  0,
				RecDiscardedCnt: 2,
			},
			callsFn: func(m *mocks.MockGeoRepoI) {},
		},
		{
			name: "check float type constraint",
			data: `ip_address,country_code,country,city,latitude,longitude,mystery_value
123.123.123.123,LI,Guyana,Port Karson,text,-163.26218895343357,1337885276
123.123.123.123,LI,Guyana,Port Karson,-78.2274228596799,text,1337885276
123.123.123.123,LI,Guyana,Port Karson,-78.2274228596799999999999999999999999999999999999,-163.26218895343357,1337885276
123.123.123.123,LI,Guyana,Port Karson,-78.2274228596799,-163.2274228596799999999999999999999999999999999999,1337885276`,
			stats: ParsingStats{
				RecAcceptedCnt:  2,
				RecDiscardedCnt: 2,
			},
			callsFn: func(m *mocks.MockGeoRepoI) {
				m.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(int64(2), nil)
			},
		},
		{
			name: "insert batch after loop break",
			data: `ip_address,country_code,country,city,latitude,longitude,mystery_value
123.123.123.123,LI,Guyana,Port Karson,-78.1245233,-163.26218895343357,1337885276
123.123.123.123,LI,Guyana,Port Karson,-78.1245233,-163.26218895343357,1337885276
123.123.123.123,LI,Guyana,Port Karson,-78.2274228596799,-78.1245233,1337885276`,
			stats: ParsingStats{
				RecAcceptedCnt:  3,
				RecDiscardedCnt: 0,
			},
			callsFn: func(m *mocks.MockGeoRepoI) {
				m.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(int64(2), nil)
				m.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(int64(1), nil)
			},
		},
		{
			name: "insert all batches before loop break",
			data: `ip_address,country_code,country,city,latitude,longitude,mystery_value
123.123.123.123,LI,Guyana,Port Karson,-78.1245233,-163.26218895343357,1337885276
123.123.123.123,LI,Guyana,Port Karson,-78.2274228596799,-78.1245233,1337885276`,
			stats: ParsingStats{
				RecAcceptedCnt:  2,
				RecDiscardedCnt: 0,
			},
			callsFn: func(m *mocks.MockGeoRepoI) {
				m.EXPECT().SaveAll(gomock.Any(), gomock.Any()).Return(int64(2), nil)
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := createFile(filePath, c.data); err != nil {
				t.Errorf("Failed to create file %q: %v\n", filePath, err)
			}

			ctrl := gomock.NewController(t)

			geoRepo := mocks.NewMockGeoRepoI(ctrl)
			c.callsFn(geoRepo)

			geoService := NewGeoService(geoRepo, validator.New(validator.WithRequiredStructEnabled()))

			stats, err := geoService.ParseCSV(context.Background(), filePath, 2)

			require.NoError(t, err)
			assert.Equal(t, c.stats.RecAcceptedCnt, stats.RecAcceptedCnt)
			assert.Equal(t, c.stats.RecDiscardedCnt, stats.RecDiscardedCnt)
		})
	}
}

func Test_GeoService_ParseCSV_Failed(t *testing.T) {
	filePath := os.TempDir() + "/Test_GeoService_ParseCSV_Failed.txt"

	defer func() {
		if dErr := os.Remove(filePath); dErr != nil {
			t.Errorf("Failed to remove file %q: %v\n", filePath, dErr)
		}
	}()

	cases := []struct {
		name   string
		data   string
		errMsg string
	}{
		{
			name: "invalid headers",
			data: `ip_address,country_code,country,city
200.106.141.15,SI,Nepal,DuBuquemouth,-84.87503094689836,7.206435933364332,7823011346
160.103.7.140,CZ,Nicaragua,New Neva,-68.31023296602508,-37.62435199624531,7301823115`,
			errMsg: "expected 7 columns, got",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if err := createFile(filePath, c.data); err != nil {
				t.Errorf("Failed to create file %q: %v\n", filePath, err)
			}

			ctrl := gomock.NewController(t)

			geoRepo := mocks.NewMockGeoRepoI(ctrl)

			geoService := NewGeoService(geoRepo, validator.New(validator.WithRequiredStructEnabled()))

			_, err := geoService.ParseCSV(context.Background(), filePath, 2)

			assert.ErrorContains(t, err, c.errMsg)
		})
	}
}

func createFile(path, data string) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}

	defer file.Close()

	if _, err = file.WriteString(data); err != nil {
		return fmt.Errorf("failed to write data: %w", err)
	}

	if err = file.Sync(); err != nil {
		return fmt.Errorf("failed to flush data: %w", err)
	}

	return nil
}
