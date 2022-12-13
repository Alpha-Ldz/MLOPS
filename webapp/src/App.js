import * as React from 'react';
import Map from 'react-map-gl';
import 'mapbox-gl/dist/mapbox-gl.css';

function App() {
  return (
    <Map
      initialViewState={{
        latitude: 48.856416650506375,
        longitude: 2.3494125832340313,
        zoom: 11
      }}
      style={{width: 800, height: 600}}
      mapStyle="mapbox://styles/mapbox/streets-v9"
      mapboxAccessToken={process.env.REACT_APP_MAPBOX_ACCESS_TOKEN}
    />
  );
}

export default App;
