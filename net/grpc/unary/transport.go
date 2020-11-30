package unary

import (
	"io/ioutil"
	"log"

	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func fromJSON(in string, pb proto.Message) error {
	err := jsonpb.UnmarshalString(in, pb)
	if err != nil {
		log.Fatalln("读取JSON时发生错误。", err.Error())
	}
	return nil
}

func toJson(pb proto.Message) string {
	marshaler := jsonpb.Marshaler{}
	str, err := marshaler.MarshalToString(pb)
	if err != nil {
		log.Fatalln("转化为Json时发生错误。", err.Error())
	}
	return str
}

func writeToFile(fileName string, pb proto.Message) error {
	dataBytes, err := proto.Marshal(pb)
	if err != nil {
		log.Fatalln("无法序列化到bytes", err.Error())
	}
	if err := ioutil.WriteFile(fileName, dataBytes, 0644); err != nil {
		log.Fatalln("无法写入到文件", err.Error())
	}
	log.Printf("成功写入文件%s。", fileName)
	return nil
}

func readFromFile(fileName string, pb proto.Message) error {
	dataBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatalln("读物文件发生错误。", err.Error())
	}
	err = proto.Unmarshal(dataBytes, pb)
	if err != nil {
		log.Fatalln("转化为struct时发生错误。", err.Error())
	}
	return nil
}
