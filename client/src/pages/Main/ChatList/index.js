import ChatItem from "./ChatItem";

import './index.css';

export default function ChatList({records}){
  return (
    <div className='chat-list'>
      {records.map((record,index)=><ChatItem key={index} record={record}/>)}
    </div>);
}