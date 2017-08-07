package parsenginx

import (
  "testing"

)

func TestDefaultParser(t *testing.T) {

  t.Log("testing default parser")

  np := NewDefaultParser()

  got := np.logFormat
  expected := `$remote_addr - $http_x_forwarded_for - $http_x_realip - [$time_local]  $scheme $http_x_forwarded_proto $x_forwarded_proto_or_scheme "$request" $status $body_bytes_sent "$http_referer" "$http_user_agent"`

  if got != expected {
    t.Errorf("expected %v, but got %v", expected, got)
  }
  got = np.reference[0]
  expected = "$remote_addr"

  if got != expected {
    t.Errorf("expected %v, but got %v", expected, got)
  }

  got = np.reference[3]
  expected = " - "

  if got != expected {
    t.Errorf("expected %v, but got %v", expected, got)
  }

}

func TestCustomParser(t *testing.T) {


  customNginxLogFormat := `$remote_addr - "$request" - $status`
  np := NewNginxParser(customNginxLogFormat)

  got := np.reference[0]
  expected := "$remote_addr"

  if got != expected {
    t.Errorf("expected %v, but got %v", expected, got)
  }

  got = np.reference[3]
  expected = "\" - "

  if got != expected {
    t.Errorf("expected %v, but got %v", expected, got)
  }

}

func TestParseLineStatus(t *testing.T) {

  np := NewDefaultParser()

  sampleLogLine := `50.112.166.232 - 50.112.166.232, 192.33.28.238, 50.112.166.232,127.0.0.1 - - - [02/Aug/2015:16:04:19 +0000]  http https,http https,http "GET /api/v1/user HTTP/1.1" 200 3350 "https://release.dollarshaveclub.com/our-products" "Mozilla/5.0 (X11; Ubuntu; Linux x86_64; rv:39.0) Gecko/20100101 Firefox/39.0"`

  got, err := np.ParseLine(sampleLogLine, "$request")
  expected := "GET /api/v1/user HTTP/1.1"
  if err != nil {
    t.Error(err)
  }

  if got != expected {
    t.Errorf("expected %v, but got %v", expected, got)
  }



}

func TestParseLineRequest(t *testing.T) {


}
