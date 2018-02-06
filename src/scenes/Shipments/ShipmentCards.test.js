import React from 'react';
import ReactDOM from 'react-dom';
import ShipmentCards from './ShipmentCards';

it('renders without crashing', () => {
  const div = document.createElement('div');
  ReactDOM.render(<ShipmentCards shipments={null} />, div);
});

//todo: test that variations of issues (null,empty,data) render as expected
