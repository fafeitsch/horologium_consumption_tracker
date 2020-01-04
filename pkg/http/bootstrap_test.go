package http

import (
	stretch "github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"reflect"
	"strconv"
	"testing"
)

func Test_convertToConfiguration(t *testing.T) {
	tests := []struct {
		name       string
		port       string
		passwords  string
		jwtKey     string
		wantErrMsg string
	}{
		{name: "Unparsable Port", port: "http-port", wantErrMsg: "the given port \"http-port\" cannot be parsed as integer"},
		{name: "Use Random Key", port: "9551"},
		{name: "Use Random Key", port: "9551", jwtKey: "password"},
		{name: "Use Random Key", port: "9551", jwtKey: "password", passwords: "./.htpasswd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := stretch.New(t)
			bindAddress := "10.20.30.40"
			args := map[string]*string{
				"port":        &tt.port,
				"jwtKey":      &tt.jwtKey,
				"bindAddress": &bindAddress,
			}
			actual, err := ConvertToConfiguration(args)
			if tt.wantErrMsg != "" {
				require.EqualError(t, err, tt.wantErrMsg, "error with correct message expected")
				return
			} else {
				require.NoError(t, err, "the testcase should not throw an error, but did.")
			}
			//Error can not be thrown because it was catched above
			wantPort, _ := strconv.Atoi(tt.port)
			assert.Equal(wantPort, actual.Port, "ports differ")
			if tt.jwtKey == "" {
				if reflect.DeepEqual(actual.JwtKey, make([]byte, len(actual.JwtKey))) {
					t.Errorf("The key consists only of zeros, which must not be the case.")
				}
			} else {
				if !reflect.DeepEqual(actual.JwtKey, []byte(tt.jwtKey)) {
					t.Errorf("Expected byte representation of \"%s\", got something else as jwtKey.", tt.jwtKey)
				}
			}
			assert.Equal(bindAddress, actual.BindAddress, "bindAddresses differ")
			assert.NotNil(actual.AccountProvider, "account provider not set")
		})
	}
}

func Test_readHtPasswd(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		want    []Account
		wantMsg string
	}{
		{name: "Valid .htpasswd", file: "../test-resources/server/valid-htpasswd", want: []Account{
			{UserName: "admin", EncryptedPassword: "$2y$11$r2rItL7fMYjnmYdoFvS96O0Xh/F0oZYxcDaoDG3j763f4DWqOfkXe"},
			{UserName: "john", EncryptedPassword: "$2y$11$l64FtbyJG2/gGdu5ecUr5.LF3YvknNL3AzxFirdnN.uo7PiIQGpHm"},
		}, wantMsg: ""},
		{name: "Invalid .htpasswd", file: "../test-resources/server/invalid-htpasswd", wantMsg: "error at line 4: Expected exactly one \":\", got 2"},
		{name: "Non-Existing .htpasswd", file: "../test-resources/server/nonexisting-htpasswd", wantMsg: "could not open passwd file: open ../test-resources/server/nonexisting-htpasswd: no such file or directory"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert := stretch.New(t)
			got, err := readHtPasswd(tt.file)
			if tt.wantMsg != "" {
				require.EqualError(t, err, tt.wantMsg, "error messages differ")
				return
			} else {
				require.NoError(t, err, "no error expected")
			}
			require.Equal(t, len(tt.want), len(got), "account lengths differ")
			for index, account := range tt.want {
				assert.Equal(account.UserName, got[index].UserName, "user name of line %d differs", index+1)
				assert.Equal(account.EncryptedPassword, got[index].EncryptedPassword, "encrypted password of line %d differs", index+1)
			}
		})
	}
}
