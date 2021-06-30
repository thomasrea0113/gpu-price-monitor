const { Cluster } = require('puppeteer-cluster');

const puppeteer = require('puppeteer-extra');

const StealthPlugin = require('puppeteer-extra-plugin-stealth');
puppeteer.use(StealthPlugin());

(async () => {
    const cluster = await Cluster.launch({
        concurrency: Cluster.CONCURRENCY_CONTEXT,
        maxConcurrency: 4,
        puppeteer,
        puppeteerOptions: {
            headless: true,
            args: [
                // TODO apparently, this is a bad idea. Should configure the docker container to use a different UID,
                // because when running as root, these arguments are required
                '--no-sandbox',
                '--disable-setuid-sandbox',

                // '--user-agent="Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36"'
            ]
        }
    });

    const contents = {};

    await cluster.task(async ({ page, data: url }) => {
        try {
            await page.goto(url, { waitUntil: 'load', timeout: 60000 });
        }
        catch (e) {
            return
        }
        contents[url] = await page.content();
    });

    process.argv.slice(2).forEach(u => cluster.queue(u));

    await cluster.idle();
    await cluster.close();

    console.log(JSON.stringify(contents))
})();