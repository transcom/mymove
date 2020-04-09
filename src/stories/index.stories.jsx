import React from 'react';
import PropTypes from 'prop-types';

import { storiesOf } from '@storybook/react';

import { action } from '@storybook/addon-actions';
import { Button } from '@trussworks/react-uswds';
import { ReactComponent as EditIcon } from 'shared/images/edit-24px.svg';
import { ReactComponent as ChevronLeft } from 'shared/icon/chevron-left.svg';
import { ReactComponent as ChevronRight } from 'shared/icon/chevron-right.svg';
import colors from '../shared/styles/colors.scss';

const filterGroup = (filter) => Object.keys(colors).filter((color) => color.indexOf(filter) === 0);

// Buttons

const ButtonGroup = ({ className }) => (
  <div className={className} style={{ padding: '20px', display: 'flex', flexWrap: 'wrap' }}>
    <Button onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <Button className="usa-button--icon" onClick={action('clicked')}>
      <span className="icon">
        <EditIcon />
      </span>
      <span>Button</span>
    </Button>
    <Button secondary onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <Button className="usa-button--small" onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <Button className="usa-button--icon usa-button--small" onClick={action('clicked')}>
      <span className="icon">
        <EditIcon />
      </span>
      <span>Button</span>
    </Button>
    <Button secondary className="usa-button--small" onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <Button secondary className="usa-button--small usa-button--icon" onClick={action('clicked')}>
      <span className="icon">
        <EditIcon />
      </span>
      <span>Button</span>
    </Button>
    <Button className="usa-button--unstyled" onClick={action('clicked')}>
      <span>Button</span>
    </Button>
    <Button className="usa-button--unstyled" onClick={action('clicked')}>
      <span className="icon">
        <EditIcon />
      </span>
      <span>Button</span>
    </Button>
  </div>
);

ButtonGroup.defaultProps = {
  className: '',
};

ButtonGroup.propTypes = {
  className: PropTypes.string,
};

storiesOf('Components|Button', module)
  .add('default', () => <ButtonGroup />)
  .add('active', () => <ButtonGroup className="active" />)
  .add('hover', () => <ButtonGroup className="hover" />)
  .add('focus', () => <ButtonGroup className="focus" />)
  .add('disabled', () => (
    <div className="disabled" style={{ padding: '20px', display: 'flex', flexWrap: 'wrap' }}>
      <Button disabled onClick={action('clicked')}>
        <span>Button</span>
      </Button>
      <Button disabled className="usa-button--icon" onClick={action('clicked')}>
        <span className="icon">
          <EditIcon />
        </span>
        <span>Button</span>
      </Button>
      <Button disabled secondary onClick={action('clicked')}>
        <span>Button</span>
      </Button>
      <Button disabled className="usa-button--small" onClick={action('clicked')}>
        <span>Button</span>
      </Button>
      <Button disabled className="usa-button--icon usa-button--small" onClick={action('clicked')}>
        <span className="icon">
          <EditIcon />
        </span>
        <span>Button</span>
      </Button>
      <Button disabled secondary className="usa-button--small" onClick={action('clicked')}>
        <span>Button</span>
      </Button>
      <Button disabled secondary className="usa-button--small usa-button--icon" onClick={action('clicked')}>
        <span className="icon">
          <EditIcon />
        </span>
        <span>Button</span>
      </Button>
      <Button disabled className="usa-button--unstyled" onClick={action('clicked')}>
        <span>Button</span>
      </Button>
      <Button disabled className="usa-button--unstyled" onClick={action('clicked')}>
        <span className="icon">
          <EditIcon />
        </span>
        <span>Button</span>
      </Button>
    </div>
  ));

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
      <a className="disabled">This link is disabled</a>
      <p>a:focus</p>
      <a className="focus">This link is focused</a>
      <p>a small</p>
      <small>
        <a href="https://materializecss.com/sass.html">USWDS blue-warm-60v 14/16</a>
      </small>
    </div>
  ));

storiesOf('Components|Containers', module).add('all', () => (
  <div id="containers" style={{ padding: '20px' }}>
    <div className="container">
      <code>
        <b>Container Default</b>
        <br />
        .container
      </code>
    </div>
    <div className="container container--gray">
      <code>
        <b>Container Gray</b>
        <br />
        .container
        <br />
        .container--gray
      </code>
    </div>
    <div className="container container--popout">
      <code>
        <b>Container Popout</b>
        <br />
        .container
        <br />
        .container--popout
      </code>
    </div>
    <div className="container container--accent--blue">
      <code>
        <b>Container Accent Blue</b>
        <br />
        .container
        <br />
        .container--accent--blue
      </code>
    </div>
    <div className="container container--accent--yellow">
      <code>
        <b>Container Accent Yellow</b>
        <br />
        .container
        <br />
        .container--accent--yellow
      </code>
    </div>
  </div>
));

// Tables

storiesOf('Components|Tables', module)
  .add('Table Elements', () => (
    <div id="sb-tables" style={{ padding: '20px' }}>
      <hr />
      <h3>Table - default</h3>
      <div className="sb-section-wrapper">
        <div className="sb-table-wrapper">
          <code>cell-bg</code>
          <table>
            <tbody>
              <tr>
                <td>Table Cell Content</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>td:hover</code>
          <table>
            <tbody>
              <tr>
                <td className="hover">Table Cell Content</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>td:locked</code>
          <table>
            <tbody>
              <tr>
                <td className="locked">Table Cell Content</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>td-numeric</code>
          <table>
            <tbody>
              <tr>
                <td className="numeric">Table Cell Content</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>th</code>
          <table>
            <thead>
              <tr>
                <th>Table Cell Content</th>
              </tr>
            </thead>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>th—sortAscending</code>
          <table>
            <thead>
              <tr>
                <th className="sortAscending">Table Cell Content</th>
              </tr>
            </thead>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>th—sortDescending</code>
          <table>
            <thead>
              <tr>
                <th className="sortDescending">Table Cell Content</th>
              </tr>
            </thead>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>th—numeric</code>
          <table>
            <thead>
              <tr>
                <th className="numeric">Table Cell Content</th>
              </tr>
            </thead>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>th—small</code>
          <table className="table--small">
            <thead>
              <tr>
                <th>Table Cell Content</th>
              </tr>
            </thead>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>th—small—numeric</code>
          <table className="table--small">
            <thead>
              <tr>
                <th className="numeric">Table Cell Content</th>
              </tr>
            </thead>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>td—filter</code>
          <table>
            <tbody>
              <tr>
                <td className="filter">
                  <input className="usa-input" id="input-type-text" name="input-type-text" type="text" />
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <hr />
      <h3>Table - stacked</h3>
      <div className="sb-section-wrapper">
        <div className="sb-table-wrapper">
          <code>td</code>
          <table className="table--stacked">
            <tbody>
              <tr>
                <td>Table Cell Content</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>th</code>
          <table className="table--stacked">
            <tbody>
              <tr>
                <th>Table Cell Content</th>
              </tr>
            </tbody>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>th: error</code>
          <table className="table--stacked">
            <tbody>
              <tr>
                <th className="error">Table Cell Content</th>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
      <hr />
      <h3>Table controls</h3>
      <div className="sb-section-wrapper">
        <div className="sb-table-wrapper">
          <code>pagination</code>
          <div className="tcontrol--pagination">
            <Button disabled className="usa-button--unstyled" onClick={action('clicked')}>
              <span className="icon">
                <ChevronLeft />
              </span>
              <span>Prev</span>
            </Button>
            <select className="usa-select" name="table-pagination">
              <option value="1">1</option>
              <option value="2">2</option>
              <option value="3">3</option>
            </select>
            <Button className="usa-button--unstyled" onClick={action('clicked')}>
              <span>Next</span>
              <span className="icon">
                <ChevronRight />
              </span>
            </Button>
          </div>
        </div>
        <div className="sb-table-wrapper">
          <code>rows per page</code>
          <div className="tcontrol--rows-per-page">
            <select className="usa-select" name="table-rows-per-page">
              <option value="1">1</option>
              <option value="2">2</option>
              <option value="3">3</option>
            </select>
            <p>rows per page</p>
          </div>
        </div>
      </div>
      <hr />
      <h3>Data points</h3>
      <div className="sb-section-wrapper">
        <div className="sb-table-wrapper">
          <code>data-point</code>
          <table className="table--data-point">
            <thead className="table--small">
              <tr>
                <th>Label</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>Table Cell Content</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>data-pair</code>
          <div className="table--data-pair">
            <table className="table--data-point">
              <thead className="table--small">
                <tr>
                  <th>Label</th>
                  <th>Label</th>
                </tr>
              </thead>
              <tbody>
                <tr>
                  <td>Table Cell Content</td>
                  <td>Table Cell Content</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </div>
    </div>
  ))
  .add('Basic Tables', () => (
    <div style={{ padding: '20px' }}>
      <p>placeholder</p>
    </div>
  ))
  .add('Table Components', () => (
    <div style={{ padding: '20px' }}>
      <p>placeholder</p>
    </div>
  ));
