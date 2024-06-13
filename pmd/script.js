document.addEventListener('DOMContentLoaded', function() {
    const contentElement = document.getElementById('content');
    const markdownUrl = '/pmd/test.md'; // 替换为你的Markdown文件链接

    fetch(markdownUrl)
        .then(response => {
            if (!response.ok) {
                throw new Error('Network response was not ok ' + response.statusText);
            }
            return response.text();
        })
        .then(markdown => {
            // 使用 Marked.js 将 Markdown 转换为 HTML
            const htmlContent = marked.parse(markdown);
            contentElement.innerHTML = htmlContent;

            // 渲染 Mermaid 图表
            mermaid.initialize({ startOnLoad: true });
            mermaid.init(undefined, contentElement.querySelectorAll('code.language-mermaid'));
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
            contentElement.textContent = 'Failed to load content.';
        });
});

