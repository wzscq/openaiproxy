import { SlackOutlined,AliwangwangFilled  } from '@ant-design/icons';
import MarkdownPreview from '@uiw/react-markdown-preview';

import './index.css';

export default function ChatItem({record}){
  const {content,role,viewLength,length}=record;
  let showContent=content;
  if(viewLength<length){
    showContent=showContent.substring(0,viewLength)+'...';
  }
  
  return (
    <div className='chat-item'>
      <div className={'chat-item-avatar '+role} style={{display:'block'}}>
        {role==='assistant'?<SlackOutlined />:<AliwangwangFilled />}
      </div>
      <div className={`chat-item-content`}>
        {role==='user'?content:<MarkdownPreview source={showContent} wrapperElement={{"data-color-mode": "dark"}} />}
      </div>
      <div className={'chat-item-avatar user'} style={{display:'none'}}>
        <AliwangwangFilled />
      </div>
    </div>
  );
}