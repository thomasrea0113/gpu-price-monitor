const puppeteer = require('puppeteer');

(async () => {
    const browser = await puppeteer.launch({
        // TODO apparently, this is a bad idea. Should configure the docker container to use a different UID
        args: ['--no-sandbox', '--disable-setuid-sandbox'],
    });
    const page = await browser.newPage();
    await page.goto('https://example.com');
    await page.screenshot({ path: 'example.png' });

    await browser.close();
})();