import { fetchEventSource,EventStreamContentType } from '@microsoft/fetch-event-source';

class RetriableError extends Error { }
class FatalError extends Error { }

//const chatProxyApi=process.env.REACT_APP_OPENAI_HOST+"/openaiproxy/openai/chat/stream/GPT4";
const chatProxyApi=process.env.REACT_APP_OPENAI_HOST+"/openai/chat/stream/GPT4";

const chatStreamCompleteProxy=(messages,callBack)=>{
  console.log("chatStreamCompleteProxy");
  fetchEventSource(chatProxyApi, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      maxTokens:2000,
      messages:messages
    }),
    async onopen(response) {
      console.log("onopen");
      if (response.ok && response.headers.get('content-type') === EventStreamContentType) {
          return; // everything's good
      } else if (response.status >= 400 && response.status < 500 && response.status !== 429) {
          // client-side errors are usually non-retriable:
          throw new FatalError();
      } else {
          throw new RetriableError();
      }
    },
    onmessage(ev) {
      console.log(ev.data);
      callBack(ev.data);
    },
    onclose() {
      console.log("Connection closed by the server");
    },
    onerror(err) {
      console.log("onerror");
      if (err instanceof FatalError) {
        throw err; // rethrow to stop the operation
      } else {
        // do nothing to automatically retry. You can also
        // return a specific retry interval here.
      }
    }
  });
}

export {
  chatStreamCompleteProxy
}