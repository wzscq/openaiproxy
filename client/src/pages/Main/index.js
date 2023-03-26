import React, { useState } from 'react';
import { SplitPane } from "react-collapse-pane"; 

import ChatInput from './ChatInput';
import ChatList from './ChatList';
import { chatCompleteProxy } from '../../utils/gptfunctions';

import './index.css';

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

const initialRecords=[
  {content:'你好，我是万能的小助理，有什么需要我帮忙处理的问题请在下方问题框中输入后点击右下角的蓝色发送按钮，我将立即帮您处理',role:'assistant'}
];  

export default function Main(){
  const [records,setRecords]=useState(initialRecords);

  const onSend=(text)=>{
    let newRecords=[...records,{content:text,role:'user'}];
    //如果newRecords中的记录数>20个，就仅保留最后的20个记录
    if(newRecords.length>20){
      newRecords=newRecords.slice(newRecords.length-20);
    }
    //将消息发送给openaiProxy
    chatCompleteProxy(newRecords).then((content)=>{
      setRecords([...newRecords,{content,role:'assistant'}]);
    });
    setRecords([...newRecords,{content:'正在处理您的请求，请稍等 ...',role:'assistant'}]);
  }

  return (
    <div className='chat-main'>
      <div className="header">ChatGPT Proxy</div>
      <div className='content'>
        <SplitPane resizerOptions={horizontalResizerOptions} initialSizes={[80,20]} split="horizontal" collapse={false}>
            <div className='chat-list-wrapper'>
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