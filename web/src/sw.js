// @ts-check

import { registerRoute, setDefaultHandler, setCatchHandler } from 'workbox-routing';
import { CacheFirst, NetworkFirst, StaleWhileRevalidate } from 'workbox-strategies';
import { skipWaiting, clientsClaim } from 'workbox-core';
import { precacheAndRoute, matchPrecache } from 'workbox-precaching';
import { ExpirationPlugin } from 'workbox-expiration';
import { RoutifyPlugin, freshCacheData } from '@roxi/routify/workbox-plugin'
const entrypointUrl = '__app.html' // entrypoint
const fallbackImage = '404.svg'
const files = self.__WB_MANIFEST // files matching globDirectory and globPattern in rollup.config.js

const externalAssetsConfig = () => ({
  cacheName: 'external',
  plugins: [
    RoutifyPlugin({
      validFor: 60 // cache is considered fresh for n seconds.
    }),
    new ExpirationPlugin({
      maxEntries: 50, // last used entries will be purged when we hit this limit
      purgeOnQuotaError: true // purge external assets on quota error
    })]
})

precacheAndRoute(files)
skipWaiting() // auto update service workers across all tabs when new release is available
clientsClaim() // take control of client without having to wait for refresh
registerRoute(isLocalPage, matchPrecache(entrypointUrl))

registerRoute(isLocalAsset, new CacheFirst())

registerRoute(hasFreshCache, new CacheFirst(externalAssetsConfig()))

setDefaultHandler(new NetworkFirst(externalAssetsConfig()));

setCatchHandler(async ({ event }) => {
  switch (event.request.destination) {
    case 'document':
      return await matchPrecache(entrypointUrl)
    case 'image':
      return await matchPrecache(fallbackImage)
    default:
      return Response.error();
  }
})

function isLocalAsset({ url, request }) { return url.host === self.location.host && request.destination != 'document' }
function isLocalPage({ url, request }) { return url.host === self.location.host && request.destination === 'document' }
function hasFreshCache(event) { return !!freshCacheData(event) }

function hasWitheringCache(event) {
  const cache = freshCacheData(event)
  if (cache) {
    const { cachedAt, validFor, validLeft, validUntil } = cache
    return validFor / 2 > validFor - validLeft
  }
}