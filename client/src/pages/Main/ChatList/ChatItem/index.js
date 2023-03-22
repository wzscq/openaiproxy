import { SlackOutlined,AliwangwangFilled  } from '@ant-design/icons';
import ReactMarkdown from 'react-markdown'

import CustomRenderers from './CustomerRenders';

import './index.css';

export default function ChatItem({record}){
  const {content,role}=record;
  return (
    <div className='chat-item'>
      <div className={'chat-item-avatar assistant'} style={{display:role==='assistant'?'block':'none'}}>
        <SlackOutlined />
      </div>
      <div className={`chat-item-content ${role}`}>
        {role==='user'?content:<ReactMarkdown renderers={CustomRenderers}>{content}</ReactMarkdown>}
      </div>
      <div className={'chat-item-avatar user'} style={{display:role==='user'?'block':'none'}}>
        <AliwangwangFilled />
      </div>
    </div>
  );
}