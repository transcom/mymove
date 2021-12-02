import React from 'react';
import { action } from '@storybook/addon-actions';
import { Button } from '@trussworks/react-uswds';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';

import QueueTable from '../components/QueueTable';
import ServiceItemsTable from '../components/Office/ServiceItemsTable/ServiceItemsTable';
import DataTable from '../components/DataTable';
import DataTableWrapper from '../components/DataTableWrapper';

const DataTableBody = (
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
  title: 'Components/Tables',
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
            <tr className="error">
              <th>Table Cell Content</th>
              <td>Table Cell Content</td>
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
              <FontAwesomeIcon icon="chevron-left" aria-labelledby="prev-label" />
            </span>
            <span id="prev-label">Prev</span>
          </Button>
          <select className="usa-select" name="table-pagination" aria-label="select page">
            <option value="1">1</option>
            <option value="2">2</option>
            <option value="3">3</option>
          </select>
          <Button className="usa-button--unstyled" onClick={action('clicked')}>
            <span id="next-label">Next</span>
            <span className="icon">
              <FontAwesomeIcon icon="chevron-right" aria-labelledby="next-label" />
            </span>
          </Button>
        </div>
      </div>
      <div className="sb-table-wrapper">
        <code>rows per page</code>
        <div className="tcontrol--rows-per-page">
          <select className="usa-select" name="table-rows-per-page" aria-label="select rows per page">
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
        <DataTable columnHeaders={['Receiving agent']} dataRow={[DataTableBody]} />
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
        <DataTableWrapper className="table--data-point-group">
          <DataTable
            columnHeaders={['Customer requested pick up date', 'Scheduled pick up date']}
            dataRow={['Thursday, 26 Mar 2020', 'Friday, 27 Mar 2020']}
          />
        </DataTableWrapper>
      </div>
    </div>
  </div>
);

export const StandardTables = () => (
  <div id="sb-tables" style={{ padding: '20px', minWidth: '1240px' }}>
    <hr />
    <h3>Queue table</h3>
    <QueueTable />
    <br />
    <hr />
    <div className="display-flex">
      <div style={{ 'margin-right': '1em' }}>
        <h3>Data point</h3>
        <DataTable columnHeaders={['Receiving agent']} dataRow={[DataTableBody]} />
      </div>
      <div style={{ 'margin-right': '1em' }}>
        <h3>Data point compact</h3>
        <DataTable
          columnHeaders={['Receiving agent']}
          dataRow={[DataTableBody]}
          custClass="table--data-point--compact"
        />
      </div>
      <div style={{ width: '40px' }} />
      <div>
        <h3>Data point group</h3>
        <DataTableWrapper className="table--data-point-group">
          <DataTable
            columnHeaders={['Customer requested pick up date', 'Scheduled pick up date']}
            dataRow={['Thursday, 26 Mar 2020', 'Friday, 27 Mar 2020']}
          />
        </DataTableWrapper>
        <div style={{ 'margin-bottom': '1em' }} />
        <DataTableWrapper className="table--data-point-group">
          <DataTable
            columnHeaders={['Authorized addresses', '']}
            dataRow={['San Antonio, TX 78234', 'Tacoma, WA 98421']}
            icon={<FontAwesomeIcon icon="arrow-right" />}
          />
          <DataTable
            columnHeaders={["Customer's addresses", '']}
            dataRow={['812 S 129th St, San Antonio, TX 78234', '441 SW Rio de la Plata Drive, Tacoma, WA 98421']}
            icon={<FontAwesomeIcon icon="arrow-right" />}
          />
        </DataTableWrapper>
      </div>
    </div>
  </div>
);

export const ServiceItemTables = () => (
  <div id="sb-tables" style={{ padding: '20px', minWidth: '1240px' }}>
    <hr />
    <h3>Service item table</h3>
    <ServiceItemsTable
      statusForTableType="SUBMITTED"
      serviceItems={[
        {
          id: 'abc12345',
          createdAt: '2020-11-22T00:00:00',
          serviceItem: 'Dom. Crating',
          code: 'DCRT',
          details: {
            description: "Here's the description",
            itemDimensions: { length: 8400, width: 2600, height: 4200 },
            crateDimensions: { length: 110000, width: 36000, height: 54000 },
          },
        },
      ]}
    />
  </div>
);
