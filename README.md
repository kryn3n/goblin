# ![goblin](assets/Goblin.png)

Goblin is a CLI tool for extracting data from ArcGIS Map/Feature Servers. It is written in Go with concurrency in mind, so it excels in extracting relatively large datasets (>1GB). Goblin was created as an alternative to GDAL's [ESRIJSON driver](https://gdal.org/en/stable/drivers/vector/esrijson.html). We found that `ogr2ogr` struggles particularly with larger datasets due to it's non-concurrent nature and the ArcGIS default record limit (usually 1-2K), which results in very long download times and/or disconnection from the server.

## Installation

### 🍺 Brew (MacOS)

The easiest way to install on MacOS is via [Homebrew](https://brew.sh/). Firstly, tap the [third-party repository](https://github.com/kryn3n/homebrew-kryn3n) with `brew tap kryn3n/kryn3n`. You can read more about third-party repositories [here](https://docs.brew.sh/Taps).

Now we can install the program as usual with `brew install kryn3n/kryn3n/goblin`.

## Usage

The `goblin` command is run with 2 required positional arguments, plus 3 optional flags. The first positional argument should be the URL to the MapServer or FeatureServer, the second should be the output filename (recommended to use `.geojsonl` as Goblin creates a [newline-delimited GeoJSON](https://en.wikipedia.org/wiki/GeoJSON#Newline-delimited_GeoJSON) file. You can read more about the format [here](https://stevage.github.io/ndgeojson/)). The optional flags can be found below.

| Flag | Description       | Required | Default |
| ---- | ----------------- | -------- | ------- |
| `-c` | Concurrency limit | No       | 100     |
| `-l` | Layer ID          | No       | -1      |
| `-m` | Max. record limit | No       | 1000    |

If no layer ID is supplied, you will be prompted for it after running your command.

# ![goblin](assets/Goblin.gif)

## Benchmarks

For this benchmark we are using the cadastre layer from Queensland's Land Parcel Property Framework server. The layer contains about 3.4 million records and the resulting file once downloaded is over 3.5GB.

Goblin downloaded the same data in less than 2 mins vs. GDAL's 105 mins. That's over 60x faster! 🚀 Results may vary depending on internet connection and the environment where the commands are being run so we encourage you to try the below commands for yourself.

### Goblin

```bash
time goblin -l 4 \
"https://spatial-gis.information.qld.gov.au/arcgis/rest/services/PlanningCadastre/LandParcelPropertyFramework/MapServer" \
test.geojsonl
________________________________________________________
Executed in  103.10 secs    fish           external
   usr time  252.44 secs    0.25 millis  252.44 secs
   sys time   23.99 secs    1.07 millis   23.99 secs
```

### GDAL ogr2ogr

```bash
time ogr2ogr test.geojsonl \
-f "GeoJSONSeq" \
"https://spatial-gis.information.qld.gov.au/arcgis/rest/services/PlanningCadastre/LandParcelPropertyFramework/MapServer/4/query?where=1=1&returnGeometry=true&outFields=*&orderByFields=objectid&f=geojson"
________________________________________________________
Executed in  105.17 mins    fish           external
   usr time  512.13 secs  715.00 micros  512.13 secs
   sys time   10.11 secs  826.00 micros   10.11 secs
```
