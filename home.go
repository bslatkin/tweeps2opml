/*
 * Copyright 2014 Brett Slatkin
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"net/http"
	"text/template"
)

var (
	homepageTemplate = template.Must(template.New("home").Parse(`
<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>Tweeps 2 OPML</title>
    <style>
        body {
            font-family: Helvetica, Arial, sans-serif;
            max-width: 600px;
        }
        h1 {
            font-size: 32px;
        }
        h2 {
            font-size: 24px;
            font-weight: normal;
        }
        .about {
            margin-top: 40px;
        }
        .submit {
            background-image: url('data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAJ4AAAAcCAYAAACQ/QaoAAAAGXRFWHRTb2Z0d2FyZQBBZG9iZSBJbWFnZVJlYWR5ccllPAAACrBJREFUeNrsW2lMVFkWPq+oKiiWKpBlEBBUbEXbHbVx32iCMZo4Gkz7w8S4jEsUbaMZpdUZp0UdjTEajTEal+gPl7jr2KDGhbjHGHVkaXVEtFXEBSwKWQrmfQdvzeNRRRUTynYm70uu79333r333HO/e5ZbKN2+fZtkfCOXNXJJlouZNGhoeZTJ5Zxc/pyYmPirJBOvk1y5HhMTExwaGko+Pj6aijS0OOx2O719+5aeP3/+Qa4m6eV/MkG6iIgITTsavAYYtM8cC5bJlwniJbdq1Yrq6uo07WjwOuBVZeKlgHhmnU6nEU/DF4EkSbgEgnga6TR8cTDxamtrNU1o+PLEc2fxzv76G6V+E+Xy/T+LS2nambtUVF6D3igpMoj+OjiBvo2w0LWiEv6mf5swTdsaPCfetedv6cdLBVQn6Sg1/g+N3heV2SjtH7lkNYWTPtiEzuhWZQWl/ZJPMX7sz+nHxDjNnWtonqt9+P4TlQeE0p+uPKOfSm00tWdcg/eHCorJGhhKOnMIST4GflZXU03llTbKr66iIMlOvxSVUnJcqNP+79y5Q69fv+b7rl27Ups2bfi+oqKCCgsLKTw8nDOhlgTOk968eUNxcXFkMpm8olh38jdnfkJeV0hISPiq5t4iFi+1fTj9XFhLUqCFfn5YQlnP3tH8Xm0oKaZeWXnWWtL5B5GkN/4nczEY5bqe6uw1ND3OQOmdgqm6urqR4lesWEE3btxo8DwlJYUyMjJ4UWbMmEFLliyh1NTUFp30zZs3afXq1bR9+3bq1KmTVxSrlv/KlSt0+fJlmjx5Mm+u5sxPyOsKly5d+q/nfvjwYcrPz2edC6xatYrfTZgw4fcjXk1NDW3oFUYLc23kExpJt6yl9MOFf1G0/hElRYfQR10AkY/eWeJM/S06mtfRQlVVVY3eHjhwgEm3fPlyGj16NJWUlPCEs7KyaPDgwdSvXz/avXs3hYSEtLibTkpK4r5bt27ttRAAFkUpPzYa5obFVI7pyVFWfHw864k9zKFDlJubS9OnT2f5m3sqoZ57QUEBy7V06VLHN6jj2MOb4ZGPvOv+Ig6QnZWspyX0t/tvyc8/gKr0viT5mkhnCiCrjy/lVhC9qDOSzuALDTbsudZOCfpKSokwOu337Nmz9PjxYxo4cCBbAKPRSO3bt6fu3btTVFQUVVZW0v79+8nPz4/fow1255w5c+jUqVMUGRlJ+/btYzfduXNnfnfs2DEeeuHChbR161Zu06NHj0Zj4/fp06dP84J++PCBtmzZQuXl5XTw4EFe4KtXr1KHDh3YBSrb5eXlsaXAt2iLzbJx40aHDPgmMzOT7t69yzIL+WHdsNHg5oqLi8nX15cCAwPp5MmT1KtXLzpy5IhjXMzfbDY3GDcsLIzlwZj37t2jR48e0ezZs3lz4tmZM2dY9oCAANYLwpddu3axLiDHs2fPHHN8//69Y+4g2Pnz5/n5ixcveLPguydPntDHjx+ZlIMGDSKbzcZzWbBgAZMWMnXs2JH0sleDJcc6CEsJ+dDGFZ9Q3r17RzqxY1wVuFqLvx+V1dRbMbhUyRRIOosc11nCSOcnWzxd49936+S4sXOgj8t++/Tpw99hoWbOnMmKwuQhNCaFyUIxuAqiQimjRo2iWbNmsSKys7N5EfAeV9ShVLzv378/7dmzhxWjHlvZNywR2q1du5a6dOnCyoXCQSh1u7Zt2/K3IDjqsDyoQy5BTNSDg4MbjAEytGvXzhHHoi6syaZNmxqMi8VTjwuvA6+BItrht0/xDNYL4z548IDf44p6Tk5OAzlBaKVc0dHRBKMDgMQgLq4AnuMe7Xfu3Mm6hIwouBf6EfrDxsFGxVya4pOQ321yYZJN7oHB0bTooZWul9Ux+Qinz5IPSbom/qBAju/GxZpd9o1YDm4GlgqKgdIB7O5169Y1MPPoA8SEVUEsAqX37duXxo4dS2r58d5isXDQfu3aNVaMKxmUY0ybNo3Gjx/P99evX+e26nYGg4EmTpzI1gsBOiwDfn+EFQPpYFkAYWUFYK0SExN588DVwdogrvJ0XGfAN+I7jAc5YDEnTZrEV9RBsEWLFrEFBHr37s1xppj7gAEDuA7ZsR4gOK4gPzYZ7l++fMkWGXKOGzeO28Ia7tixg8cSgF7mzZvHawNv5Q465SScFXRkK7dSgrGavjPr+LjELeRvkiwSRRhc9/3p0ycaOXIkTwALuWzZMs7O4L5OnDjhWDhc8T0WOjY2lt0C2iozMuVC4TnegyTqRRJF2be4B1FBUhRYBmftoAu4Qk6q5MUCSZAgYJHhTrHAuO/WrVuDMZBYoS2Ae5TmjKssyk0jnsHqjRkzhskMuXCFXEJOEDAtLY1/qFfKBYKIOnQG4uEq3uMe4QGAdRoyZAgX3IsEUbTHWlqtVsdGdzcHjw6Qo/wNZK+ULd6bGpL8TPVWryne2aspIyHIoWyn2fLnTO7o0aNs1keMGMG7bOrUqUwutWUCKbGwcBEgF6ykM8slxlQ+8yRIVpKwqXY9e/bkK+IqLHBMTAwNGzaMLl68yBkjFlidTDXVr6fjurPYcOHA3r17HURAPCrqsLJqS+pOR3gG9yssM+JxJeCdhJUHaZuTjLi1eCjoNKNHBOUMD6f0WD3FGGtdWz45qVgW70ttjU33CYIBmzdv5mAfQem2bdv4GdySUJLY2XBHsHrI5hCPrVy5ssF7MWn1zlJaBrXFc/Wds75EQXiA+BEkA5AQYVFFHRbRWb+iT1gfELa547r7Bu4cuHXrFm9qWB5sCNSFm3XWXtThTuFthFzY5JAV8SM2/fHjx1n/KEjc4MKVltsdh9Tr4ja5UAbkZXIGWG2zUmlFFf805iSjoCmtdTQpyuDYAa4KiAQFXbhwgc+y0tPTWUnz58+n5OTkRhYBFnHNmjVMSgSxIJ8nFqSpANeZxXPXFoobPny4w2pD+chMBdTxnWiH+BQLiMSkqKio2eO6soyiwI3C2gKQB8/g8vnnSnmjIIt2Ns7QoUM5PIB1RLaJZ1gDEGzDhg28jjhvhc4XL17M5enTp3xygIzdU7nVc8BfINchu3GFn3Ie0W+2anpRYaeXdj3/QqELCiHJ6FefZHyGWbLT3zubaGiwxCR1B5wTQRlYOGGuEcPBoqA9YrSgoCC+R6xx7tw5un//PgexyHpxP3fuXJoyZQqTGN+iDZQnEgFleyWgMH9/fw6SAfV36r7UEO0hO/qAzMhkARxXQLnq8XH0gPmK+SEkaO64ym8wrvpQHv3jWEr0iWMbuEqQp6ysrNHc0R7foI71EM9Qx3cAZBGyiwRQ/LqEGFHdnyfA6QUTTxxEOgMGPFpYSueLKyn7XS2f40l6Q/0Riqzg78wSjY800PcRRtJV2pweFrsjICbmLk5AkIvdhveIZ3BkgLbr16/nQ1oN3odI2Jobz6mBTJmJh3OlpgBWYyehqAGWg/nKDMkrP7HI5MSEEUvBWiADxFkgxsdzDf87ePXqlWdZrUj3f0+AYIhjEJOIYFaZzmv4P/vrlK8JkNPTOELD10+8MnlBzZIkadrQ4HV89lBWEO+cnAH9UWQxGjR4E58z+CwQL6O8vHyEzMRgpNaa5dPgLUuHBNRms+E/dGeAeHlySZIfZMrle/k+SFOTBi8Ah6bZclmamJiY/28BBgBzgDaRX4k8wAAAAABJRU5ErkJggg==');
            width: 158px;
            height: 28px;
            border: 0;
        }
    </style>
    {{if .Globals.AnalyticsId}}
    <script>
      (function(i,s,o,g,r,a,m){i['GoogleAnalyticsObject']=r;i[r]=i[r]||function(){
      (i[r].q=i[r].q||[]).push(arguments)},i[r].l=1*new Date();a=s.createElement(o),
      m=s.getElementsByTagName(o)[0];a.async=1;a.src=g;m.parentNode.insertBefore(a,m)
      })(window,document,'script','//www.google-analytics.com/analytics.js','ga');

      ga('create', '{{.Globals.AnalyticsId}}', 'auto');
      ga('send', 'pageview');
    </script>
    {{end}}
</head>
<body>
    <div class="header">
        <h1>Tweeps 2 OPML</h1>
        <h2>Get the RSS feeds of all of my Twitter friends</h2>
        <p>
            This tool will find all of the RSS feeds that your Twitter friends provide and put them all into a single OPML file that you can then import into NewsBlur, The Old Reader, Feedly, Digg, etc.
        </p>
        <p>
            Start by clicking this button:
        </p>
    </div>
    <form action="/authorize" method="POST">
        <input class="submit" type="submit" value="" />
        <input type="hidden" name="continue" value="/list" />
        <input type="hidden" name="callback" value="/oauth_callback" />
    </form>
    <div class="about">
        <a href="https://github.com/bslatkin/tweeps2opml">About & source code</a>
    </div>
</body>
</html>
`))
)

func homepageHandler(c *Context, w http.ResponseWriter, r *http.Request) {
	if r.RequestURI != "/" {
		http.Error(w, "Unknown URL", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	homepageTemplate.Execute(w, Params{Globals: Globals})
}
