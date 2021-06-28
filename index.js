const puppeteer = require('puppeteer-extra');

const StealthPlugin = require('puppeteer-extra-plugin-stealth');
puppeteer.use(StealthPlugin());

(async () => {
    const browser = await puppeteer.launch({
        // headless: true,
        args: [
            // TODO apparently, this is a bad idea. Should configure the docker container to use a different UID,
            // because when running as root, these arguments are required
            '--no-sandbox',
            '--disable-setuid-sandbox',

            '--user-agent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"'],
    });
    const page = await browser.newPage();
    await page.goto(process.argv[2])
    const content = await page.content();
    await browser.close()
    console.log(content);
})();