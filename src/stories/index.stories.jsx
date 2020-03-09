import React from 'react';
import PropTypes from 'prop-types';

import { storiesOf } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { linkTo } from '@storybook/addon-links';
import colors from '../shared/styles/colors.scss';

import { Button, Welcome } from '@storybook/react/demo';

const filterGroup = filter => Object.keys(colors).filter(color => color.indexOf(filter) === 0);

storiesOf('Welcome', module).add('to Storybook', () => <Welcome showApp={linkTo('Button')} />);

storiesOf('Button', module)
  .add('with text', () => <Button onClick={action('clicked')}>Hello Button</Button>)
  .add('with some emoji', () => (
    <Button onClick={action('clicked')}>
      <span role="img" aria-label="so cool">
        😀 😎 👍 💯
      </span>
    </Button>
  ));

storiesOf('Global|Colors', module).add('all', () => (
  <div style={{ padding: '20px' }}>
    <>
      <h3>Brand Colors</h3>
      <ColorGroup group={filterGroup('brand')} />
    </>
    <>
      <h3>Background Colors</h3>
      <ColorGroup group={filterGroup('background')} />
    </>
    <>
      <h3>Base Colors</h3>
      <ColorGroup group={filterGroup('base')} />
    </>
    <>
      <h3>Alert Colors</h3>
      <ColorGroup group={filterGroup('alert')} />
    </>
    <>
      <h3>Accent Colors</h3>
      <ColorGroup group={filterGroup('accent')} />
    </>
  </div>
));

// Convert the color key to the color variable name.
const colorVariable = color => {
  const array = color.split('-')[1].split(/(?=[A-Z])/);
  return `$color-${array.join('-').toLowerCase()}`;
};

// Convert the color key to the color proper name.
const colorName = color => {
  const array = color.split('-')[1].split(/(?=[A-Z])/);
  return `${array.join(' ').toLowerCase()}`;
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
        backgroundColor: colors[color],
        display: 'block',
        height: '4em',
        marginBottom: '0.3em',
        borderRadius: '5px',
        border: '1px solid lightgray',
      }}
    />
    <span>{colorName(color)}</span> <br />
    <span>{colorVariable(color)}</span> <br />
    <span>{colors[color]}</span> <br />
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
    {group.map(color => {
      return <Color color={color} key={color} />;
    })}
  </ul>
);

ColorGroup.propTypes = {
  group: PropTypes.array.isRequired,
};
