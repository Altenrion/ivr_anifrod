package main

import (
	"fmt"
	"github.com/bclicn/color"
)

type NatsWriter struct{
	err error
}

func (p *NatsWriter) Write(w []byte) (n int, err error) {
	if err != nil{
		return 0, err
	}

	hash := "user.ivr.action"
	data := w

	if NatsErr != nil {
		fmt.Printf(color.Red(" error :[%s], row :[%s] \n"), NatsErr, w)
	}else{
		fmt.Printf(color.LightGreen("[%s] \n"), w)
		errOnPublish := NatsClient.Publish(hash, data)
		if errOnPublish != nil {
			p.err = errOnPublish
			return 0, errOnPublish
		}
	}
	return len(w), nil
}

