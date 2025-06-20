import React from 'react';
import PropTypes from 'prop-types';

import colors from 'styles/colors.module.scss';

const filterGroup = (filter) => Object.keys(colors).filter((color) => color.indexOf(filter) === 0);

// Colors

export default {
  title: 'Global/Colors',
};

export const BrandColors = () => (
  <div style={{ padding: '20px' }}>
    <h3>Brand Colors</h3>
    <ColorGroup group={filterGroup('brand')} />
  </div>
);
export const BackgroundColors = () => (
  <div style={{ padding: '20px' }}>
    <h3>Background Colors</h3>
    <ColorGroup group={filterGroup('background')} />
  </div>
);
export const BaseColors = () => (
  <div style={{ padding: '20px' }}>
    <h3>Base Colors</h3>
    <ColorGroup group={filterGroup('base')} />
  </div>
);
export const AlertColors = () => (
  <div style={{ padding: '20px' }}>
    <h3>Alert Colors</h3>
    <ColorGroup group={filterGroup('alert')} />
  </div>
);
export const PrimaryColors = () => (
  <div style={{ padding: '20px' }}>
    <h3>Primary Colors</h3>
    <ColorGroup group={filterGroup('primary')} />
  </div>
);
export const AccentColors = () => (
  <div style={{ padding: '20px' }}>
    <h3>Accent Colors</h3>
    <ColorGroup group={filterGroup('accent')} />
  </div>
);

// Convert the color key to the color variable name.
const colorVariable = (color) => {
  const array = color.split('-')[1].split(/(?=[A-Z])/);
  return `$${array.join('-').toLowerCase()}`;
};

// Convert the color key to the color proper name.
const colorName = (color) => {
  const array = color.split('-')[1].split(/(?=[A-Z])/);
  return `${array.join(' ').toLowerCase()}`;
};

const colorsHelper = (color) => {
  if (Object.keys(colors).includes(color)) {
    // The use of colors[color] triggers a security warning from our eslint security plugin.
    // However, since we verify inputs against imported colors and this function is not used where
    // users input color we are diabling the warning.
    return colors[color];
  }
  return colors.base;
};

// A component for displaying individual color swatches.
const Color = ({ color }) => (
  <li
    style={{
      borderRadius: '5px',
      border: '1px solid lightgray',
      padding: '5px',
    }}
  >
    <span
      style={{
        backgroundColor: colorsHelper(color),
        display: 'block',
        height: '4em',
        marginBottom: '0.3em',
        borderRadius: '5px',
        border: '1px solid lightgray',
      }}
    />
    <p
      style={{
        fontSize: '13px',
      }}
    >
      <span style={{ 'text-transform': 'capitalize' }}>
        <b>{colorName(color)}</b>
      </span>
      <br />
      <code>{colorVariable(color)}</code>
      <br />
      <code>{colorsHelper(color)}</code>
      <br />
    </p>
  </li>
);

Color.propTypes = {
  color: PropTypes.string.isRequired,
};

// A component for displaying a group of colors.
const ColorGroup = ({ group }) => (
  <ul
    style={{
      display: 'grid',
      gridTemplateColumns: 'repeat(auto-fill, minmax(120px, 175px))',
      gridGap: '20px',
      marginBottom: '40px',
      listStyle: 'none',
      padding: '0px',
    }}
  >
    {group.map((color) => {
      return <Color color={color} key={color} />;
    })}
  </ul>
);

ColorGroup.propTypes = {
  group: PropTypes.arrayOf.isRequired,
};
