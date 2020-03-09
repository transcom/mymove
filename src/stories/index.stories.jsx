import React from 'react';
import PropTypes from 'prop-types';

import { storiesOf } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { linkTo } from '@storybook/addon-links';
import colors from '../shared/styles/colors.scss';

import { Button, Welcome } from '@storybook/react/demo';

const filterGroup = (filter) =>
  Object.keys(colors).filter((color) => color.indexOf(filter) === 0);

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

  const colors = () => {
    return (
      <ul>
        {Object.keys(colors).map((color) => (
          <li>
            <span
              style={
                backgroundColor: colors[color],
                display: 'block',
                height: '4em',
                marginBottom: '0.3em',
                borderRadius: '5px',
                border: '1px solid lightgray'
              }
            />
          <span>{color}</span><br /> // color name
          <span>{colors[color]}</span> <br /> // hex value
        </li>
        )
      )
    )
  }
