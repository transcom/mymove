import React from 'react';
import { action } from '@storybook/addon-actions';
import { Button } from '@trussworks/react-uswds';

import QueueTable from '../components/QueueTable';
import ServiceItemTable from '../components/ServiceItemTable';
import ServiceItemTableHasImg from '../components/ServiceItemTableHasImg';
import DataPoint from '../components/DataPoint';
import DataPointGroup from '../components/DataPointGroup';

import { ReactComponent as ChevronRight } from 'shared/icon/chevron-right.svg';
import { ReactComponent as ChevronLeft } from 'shared/icon/chevron-left.svg';
import { ReactComponent as ArrowRight } from 'shared/icon/arrow-right.svg';

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

export default {
  title: 'Components|Tables',
  parameters: {
    abstract: {
      url: 'https://share.goabstract.com/0a7c55ae-8268-4298-be70-4cf0117c6034?mode=design',
    },
  },
};

export const TableElements = () => (
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
        <DataPoint columnHeaders={['Receiving agent']} dataRow={[dataPointBody]} />
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
        <DataPointGroup>
          <DataPoint
            columnHeaders={['Customer requested pick up date', 'Scheduled pick up date']}
            dataRow={['Thursday, 26 Mar 2020', 'Friday, 27 Mar 2020']}
          />
        </DataPointGroup>
      </div>
    </div>
  </div>
);

export const StandardTables = () => (
  <div id="sb-tables" style={{ padding: '20px' }}>
    <hr />
    <h3>Queue table</h3>
    <QueueTable />
    <br />
    <hr />
    <div className="display-flex">
      <div style={{ 'margin-right': '1em' }}>
        <h3>Data point</h3>
        <DataPoint columnHeaders={['Receiving agent']} dataRow={[dataPointBody]} />
      </div>
      <div style={{ 'margin-right': '1em' }}>
        <h3>Data point compact</h3>
        <DataPoint
          columnHeaders={['Receiving agent']}
          dataRow={[dataPointBody]}
          custClass="table--data-point--compact"
        />
      </div>
      <div style={{ width: '40px' }} />
      <div>
        <h3>Data point group</h3>
        <DataPointGroup>
          <DataPoint
            columnHeaders={['Customer requested pick up date', 'Scheduled pick up date']}
            dataRow={['Thursday, 26 Mar 2020', 'Friday, 27 Mar 2020']}
          />
        </DataPointGroup>
        <div style={{ 'margin-bottom': '1em' }} />
        <DataPointGroup>
          <DataPoint
            columnHeaders={['Authorized addresses', '']}
            dataRow={['San Antonio, TX 78234', 'Tacoma, WA 98421']}
            Icon={ArrowRight}
          />
          <DataPoint
            columnHeaders={["Customer's addresses", '']}
            dataRow={['812 S 129th St, San Antonio, TX 78234', '441 SW Rio de la Plata Drive, Tacoma, WA 98421']}
            Icon={ArrowRight}
          />
        </DataPointGroup>
      </div>
    </div>
  </div>
);

export const ServiceItemTables = () => (
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
          submittedAt: '2020-11-22',
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
);
