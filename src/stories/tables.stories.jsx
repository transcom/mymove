import React from 'react';
import { storiesOf } from '@storybook/react';
import { action } from '@storybook/addon-actions';
import { Button } from '@trussworks/react-uswds';

import QueueTable from '../components/QueueTable';
import ServiceItemTable from '../components/ServiceItemTable';
import ServiceItemTableHasImg from '../components/ServiceItemTableHasImg';
import DataPoint from '../components/DataPoint';
import DataPair from '../components/DataPair';

import { ReactComponent as ChevronLeft } from 'shared/icon/chevron-left.svg';
import { ReactComponent as ChevronRight } from 'shared/icon/chevron-right.svg';

const dataPointBody = (
  <>
    Dorothy Lagomarsino
    <br />
    +1 999-999-9999
    <br />
    dorothyl@email.com
  </>
);

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
                <td>Table Cell Content</td>
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
                <th>Table Cell Content</th>
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
              <tr className="filter">
                <td>Table Cell Content</td>
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
          <code>data point</code>
          <br />
          <br />
          <DataPoint header="Receiving agent" body={dataPointBody} />
          <br />
          <br />
          <code>data point compact</code>
          <br />
          <br />
          <table className="table--data-point table--data-point--compact">
            <thead className="table--small">
              <tr>
                <th>Dorothy Lagomarsino</th>
              </tr>
            </thead>
            <tbody>
              <tr>
                <td>+1 999-999-9999</td>
              </tr>
            </tbody>
          </table>
        </div>
        <div className="sb-table-wrapper">
          <code>data-pair</code>
          <DataPair />
        </div>
      </div>
    </div>
  ))
  .add('Standard Tables', () => (
    <div id="sb-tables" style={{ padding: '20px' }}>
      <hr />
      <h3>Queue table</h3>
      <QueueTable />
      <br />
      <hr />
      <div className="display-flex">
        <div>
          <h3>Data point</h3>
          <DataPoint header="Receiving agent" body={dataPointBody} />
        </div>
        <div>
          <h3>Data point compact</h3>
          <DataPoint header="Receiving agent" body={dataPointBody} custClass="table--data-point--compact" />
        </div>
        <div style={{ width: '40px' }} />
        <div>
          <h3>Data pair</h3>
          <DataPair />
        </div>
      </div>
    </div>
  ))
  .add('Service Item Tables', () => (
    <div id="sb-tables" style={{ padding: '20px' }}>
      <hr />
      <h3>Service item table</h3>
      <ServiceItemTable />
      <br />
      <hr />
      <h3>Service item table with images and buttons</h3>
      <ServiceItemTableHasImg
        serviceItems={[
          {
            id: 'abc12345',
            dateRequested: '22 Nov 2020',
            serviceItem: 'Dom. Crating',
            code: 'DCRT',
            details: {
              text: {
                Description: "Here's the description",
                'Item dimensions': '84"x26"x42"',
                'Crate dimensions': '110"x36"x54"',
              },
              imgURL: 'https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg',
            },
          },
        ]}
      />
    </div>
  ));
