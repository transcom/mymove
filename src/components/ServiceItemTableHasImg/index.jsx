import React from 'react';
import { Button } from '@trussworks/react-uswds';
import { ReactComponent as Check } from '../../shared/icon/check.svg';
import { ReactComponent as Ex } from '../../shared/icon/ex.svg';

const ServiceItemTableHasImg = () => (
  <div className="table--service-item table--service-item--hasimg">
    <table>
      <col style={{ width: '120px' }} />
      <col style={{ width: '170px' }} />
      <col style={{ width: '100px' }} />
      <col style={{ width: '350px' }} />
      <col />
      <thead className="table--small">
        <tr>
          <th>Date requested</th>
          <th>Service item</th>
          <th>Code</th>
          <th>Details</th>
          <th>&nbsp;</th>
        </tr>
      </thead>
      <tbody>
        <tr style={{ height: '80px' }}>
          <td style={{ paddingTop: '19px', verticalAlign: 'top' }}>20 Nov 2019</td>
          <td style={{ paddingTop: '19px', verticalAlign: 'top' }}>Domestic crating</td>
          <td style={{ paddingTop: '19px', verticalAlign: 'top' }}>DCRT</td>
          <td style={{ verticalAlign: 'top' }}>
            <div className="display-flex" style={{ alignItems: 'center' }}>
              <div
                className="si-thumbnail"
                style={{
                  width: '56px',
                  height: '42px',
                  backgroundImage: 'url("https://live.staticflickr.com/4735/24289917967_27840ed1af_b.jpg")',
                }}
                aria-labelledby="si-thumbnail--caption"
              />
              <small id="si-thumbnail--caption">grandfather clock 7ft x 2ft x 3.5ft </small>
            </div>
          </td>
          <td>
            <div className="display-flex">
              <Button className="usa-button--icon usa-button--small">
                <span className="icon">
                  <Check />
                </span>
                <span>Accept</span>
              </Button>
              <Button secondary className="usa-button--small usa-button--icon">
                <span className="icon">
                  <Ex />
                </span>
                <span>Reject</span>
              </Button>
            </div>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
);

export default ServiceItemTableHasImg;
