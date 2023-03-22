import {Input,Button} from 'antd';
import {useState} from 'react';
import { SendOutlined } from '@ant-design/icons';

const {TextArea} = Input;

const TextAreaStyles={
  height: "calc(100% - 10px)",
  resize: 'none',
  width:'calc(100% - 60px)',
  margin:5
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
    <TextArea style={TextAreaStyles} value={text} onChange={(e)=>setText(e.target.value)} />
    <Button disabled={text.length>0?false:true} type="primary" style={{float:'right',margin:"5px 5px 5px 0px",height:"calc(100% - 10px)",width:45}} icon={<SendOutlined />} onClick={sendMessage}/>
    </>
  );
}