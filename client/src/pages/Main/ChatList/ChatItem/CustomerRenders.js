import { copy } from 'copy-to-clipboard';
import { highlight, languages } from 'prismjs/components/prism-core';

const CustomRenderers = {
  code: ({ language, value }) => {
    const className = `language-${language}`;
    const highlightedCode = highlight(value, languages[language]);
    return (
      <div>
        <pre className={className}>
          <code
            className={className}
            dangerouslySetInnerHTML={{ __html: highlightedCode }}
          />
        </pre>
        <button onClick={() => copy(value)}>Copy</button>
      </div>
    );
  },
};

export default CustomRenderers;