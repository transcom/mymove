import React from 'react';
import PropTypes from 'prop-types';
import { storiesOf } from '@storybook/react';
import { Tag } from '@trussworks/react-uswds';

import colors from '../shared/styles/colors.scss';

import { ReactComponent as AlertIcon } from 'shared/icon/alert.svg';

const filterGroup = (filter) => Object.keys(colors).filter((color) => color.indexOf(filter) === 0);

// Colors

storiesOf('Global|Colors', module).add('all', () => (
  <div style={{ padding: '20px' }}>
    <h3>Brand Colors</h3>
    <ColorGroup group={filterGroup('brand')} />
    <h3>Background Colors</h3>
    <ColorGroup group={filterGroup('background')} />
    <h3>Base Colors</h3>
    <ColorGroup group={filterGroup('base')} />
    <h3>Alert Colors</h3>
    <ColorGroup group={filterGroup('alert')} />
    <h3>Primary Colors</h3>
    <ColorGroup group={filterGroup('primary')} />
    <h3>Accent Colors</h3>
    <ColorGroup group={filterGroup('accent')} />
  </div>
));

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
    // eslint-disable-next-line security/detect-object-injection
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

// Typography
storiesOf('Global|Typography', module)
  .add('Headers', () => (
    <div style={{ padding: '20px' }}>
      <p>h1</p>
      <h1>Public Sans 40/48</h1>
      <p>h2</p>
      <h2>Public Sans 28/34</h2>
      <p>h3</p>
      <h3>Public Sans 22/26</h3>
      <p>h4</p>
      <h4>Public Sans 17/20</h4>
      <p>h5</p>
      <h5>Public Sans 15/21</h5>
      <p>h6</p>
      <h6>Public Sans 13/18</h6>
    </div>
  ))
  .add('Text', () => (
    <div style={{ padding: '20px' }}>
      <p>p</p>
      <p>
        Public Sans 15/23 Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut
        labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut
        aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu
        fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit
        anim id est laborum.
      </p>
      <p>p small</p>
      <p>
        <small>
          Public Sans 13/18 Faucibus in ornare quam viverra orci sagittis eu volutpat odio. Felis imperdiet proin
          fermentum leo vel orci. Egestas sed sed risus pretium quam vulputate. Consectetur libero id faucibus nisl.
          Ipsum dolor sit amet consectetur adipiscing elit. Id aliquet lectus proin nibh nisl condimentum id venenatis
          a. Pellentesque pulvinar pellentesque habitant morbi tristique senectus. Mattis vulputate enim nulla aliquet
          porttitor lacus luctus accumsan.
        </small>
      </p>
      <p>
        Bulleted list
        <ul>
          <li>Public Sans 15px / 23px</li>
          <li>USWDS bg-base-darkest</li>
        </ul>
      </p>
    </div>
  ))
  .add('Links', () => (
    <div style={{ padding: '20px' }}>
      <p>a</p>
      <a href="https://materializecss.com/sass.html">USWDS blue-warm-60v</a>
      <p>a:hover</p>
      <a className="hover" href="https://materializecss.com/sass.html">
        USWDS blue-warm-60v
      </a>
      <p>a:visted</p>
      <a className="visited" href="#">
        USWDS bg-violet-warm-60
      </a>
      <p>a:disabled</p>
      <a href="#" className="disabled">
        This link is disabled
      </a>
      <p>a:focus</p>
      <a href="#" className="focus">
        This link is focused
      </a>
      <p>a small</p>
      <small>
        <a href="https://materializecss.com/sass.html">USWDS blue-warm-60v 14/16</a>
      </small>
    </div>
  ));

// Tags

storiesOf('Components|Tags', module).add('all', () => (
  <div id="tags" style={{ padding: '20px' }}>
    <hr />
    <h3>Tags</h3>
    <Tag>New</Tag>
    <Tag className="usa-tag--green">Authorized</Tag>
    <Tag className="usa-tag--red">Rejected</Tag>
    <Tag className="usa-tag--yellow">Pending</Tag>
    <Tag className="usa-tag--alert">
      <AlertIcon />
    </Tag>
    <Tag className="usa-tag--teal">INTL</Tag>
    <Tag>3</Tag>
    <Tag className="usa-tag--cyan usa-tag--large">#ABC123K</Tag>
  </div>
));
