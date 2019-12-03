import React from 'react';
import ReactDOM from 'react-dom';
import RangeSlider from '.';

it('renders without crashing', () => {
  const onChangeMock = jest.fn();
  const div = document.createElement('div');
  ReactDOM.render(
    <RangeSlider defaultValue={0} min={0} max={100} step={50} id="tester" onChange={onChangeMock} />,
    div,
  );
});
