import axios from "axios";

//chatproxy
const chatProxyApi="/openaiproxy/openai/v1/chat/completions/GPT3Dot5Turbo"
const chatCompleteProxy=async (messages)=>{
  
  const reponse= await axios({
    url:chatProxyApi,
    timeout:300000,
    method:"post",
    headers: {
      'Content-Type': 'application/json'
    },
    data:{
      maxTokens:2000,
      messages:messages
    }});

  if(reponse.data?.error===true){
    return reponse.data?.message;
  }
  
  return reponse.data?.result?.choices[0]?.message?.content;
}

export {
  chatCompleteProxy
}