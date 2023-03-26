import {Input,Button} from 'antd';
import {useState} from 'react';
import { SendOutlined } from '@ant-design/icons';

const {TextArea} = Input;

const TextAreaStyles={
  height: "calc(100% - 10px)",
  resize: 'none',
  width:'calc(100% - 95px)',
  margin:5,
  marginLeft:40,
  backgroundColor:'#0d1117',
  color:'white'
}

export default function ChatInput({onSend}){
  const [text,setText]=useState('');

  const sendMessage=()=>{
    if(text.length>0){
      onSend(text);
      setText('');
    }
  }

  return (
    <>
    <TextArea placeholder={"在这里输入您的问题，然后点击右侧按钮发送"} style={TextAreaStyles} value={text} onChange={(e)=>setText(e.target.value)} />
    <Button type="primary" style={{float:'right',margin:"5px 5px 5px 0px",height:"calc(100% - 10px)",width:45}} icon={<SendOutlined />} onClick={sendMessage}/>
    </>
  );
}