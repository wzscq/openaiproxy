import React, { useEffect, useRef, useState } from 'react';
import { SplitPane } from "react-collapse-pane"; 
import {Input,Select} from "antd";

import ChatInput from './ChatInput';
import ChatList from './ChatList';
import { chatStreamCompleteProxy } from '../../utils/gptStream';

import './index.css';

const Option=Select.Option;

const InputStyles={
  backgroundColor:'#0d1117',
  color:'white'
}

const horizontalResizerOptions={
  css: {
    height: '1px',
    background: 'rgba(0, 0, 0, 0.1)',
  },
  hoverCss: {
    height: '2px',
    background: 'rgba(0, 0, 0, 0.1)',
  },
  grabberSize: '2px',
}

const initContent='你好，有什么需要我帮忙处理的问题请在下方问题框中输入后点击右下角的蓝色发送按钮，我将立即帮您处理';

const initialRecords=[
  {content:initContent,role:'assistant',length:initContent.length,viewLength:initContent.length}
];  

var gContent="";

export default function Main(){
  const [records,setRecords]=useState(initialRecords);
  const [account,setAccount]=useState(localStorage.getItem('account',''));
  const [password,setPassword]=useState(localStorage.getItem('password',''));
  const [model,setModel]=useState(localStorage.getItem('model','GPT4'));

  const refList = useRef();

  useEffect(()=>{
    localStorage.setItem('account',account);
    localStorage.setItem('password',password);
    localStorage.setItem('model',model);
  },[account,password,model]);

  const onSend=(text)=>{
    console.log('onSend');
    let newRecords=[...records,{content:text,role:'user',viewLength:text.length,length:text.length}];
    //如果newRecords中的记录数>20个，就仅保留最后的20个记录
    if(newRecords.length>20){
      newRecords=newRecords.slice(newRecords.length-20);
    }
    gContent="";
    //将消息发送给openaiProxy
    chatStreamCompleteProxy(newRecords,account,password,model,(text)=>{
      //br 替换回回车
      text=text.replaceAll(/<br\/>/g,'\n');
      gContent+=text;
      console.log(gContent);
      setRecords([...newRecords,{content:gContent,role:'assistant',viewLength:gContent.length,length:gContent.length}]);
    });
    setRecords([...newRecords,{content:'正在处理您的请求，请稍等 ...',role:'assistant'}]);
  }

  useEffect(()=>{
    const updateViewLength=(index,viewLength)=>{
      //console.log('updateViewLength',viewLength);
      const newRecords=[...records];
      newRecords[index].viewLength=viewLength;
      setRecords(newRecords);
    }

    records.forEach((record,index)=>{
      //console.log('updateViewLength',record,record.role==='assistant',record.viewLength,record.length,record.viewLength<record.length);
      if(record.role==='assistant' && record.viewLength<record.length){
        let viewLength=record.viewLength+Math.floor(Math.random() * 10);
        if(viewLength>record.length){
          viewLength=record.length;
        }
        setTimeout(()=>updateViewLength(index,viewLength),Math.floor(Math.random() * 500));
      }
    });
  },[records]);

  useEffect(()=>{
    refList.current.scrollTop = refList.current.scrollHeight;
  },[records]);

  return (
    <div className='chat-main'>
      <div className="header">
        <div className='title'>GPT Proxy</div>
        <div className='account'><Input style={InputStyles}  placeholder='请输入账号' value={account} onChange={(e)=>setAccount(e.target.value)} /></div>
        <div className='password'><Input.Password style={InputStyles}  placeholder='请输入密码' value={password} onChange={(e)=>setPassword(e.target.value)} /></div>
        <div className='model'>
          <Select value={model} onChange={(value)=>setModel(value)}>
            <Option value="GPT3.5">GPT3.5</Option>
            <Option value="GPT4">GPT4</Option>
          </Select>
        </div>
      </div>
      <div className='content'>
        <SplitPane resizerOptions={horizontalResizerOptions} initialSizes={[80,20]} split="horizontal" collapse={false}>
            <div ref={refList} className='chat-list-wrapper'>
              <ChatList records={records}/>
            </div>
            <div className='chat-input-wrapper'>
              <ChatInput onSend={onSend}/>
            </div>
        </SplitPane>
      </div>
    </div>
  );
}