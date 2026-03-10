const CACHE_NAME = "temple-app-v1";
const urlsToCache = [
    "/",
    "/static/css/global.css",
    "/static/images/temple-bg.jpg"
];

self.addEventListener("install", event => {
    event.waitUntil(
        caches.open(CACHE_NAME).then(cache => {
            return cache.addAll(urlsToCache);
        })
    );
});

self.addEventListener("fetch", event => {
    event.respondWith(
        caches.match(event.request).then(response => {
            return response || fetch(event.request);
        })
    );
});