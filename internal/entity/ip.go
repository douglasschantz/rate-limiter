package entity

import "errors"

type Ip struct {
	IP string
}

type Count struct {
	Count int
}

func NewIp(ip string) (*Ip, error) {
	respIp := &Ip{
		IP: ip,
	}
	err := respIp.IsValid()
	if err != nil {
		return nil, err
	}
	return respIp, nil
}
func (i *Ip) IsValid() error {
	if i.IP == "" {
		return errors.New("invalid ip")
	}
	return nil
}

func (o *Count) CheckLimiterIp() bool {
	return o.Count <= 10
}
