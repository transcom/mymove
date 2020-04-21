import React from 'react';
import PropTypes from 'prop-types';
import { Button } from '@trussworks/react-uswds';
import { ReactComponent as Check } from '../../shared/icon/check.svg';
import { ReactComponent as Ex } from '../../shared/icon/ex.svg';

const ServiceItemTableHasImg = ({ serviceItems }) => {
  const tableRows = serviceItems.map(({ id, dateRequested, serviceItem, code, details }) => {
    let detailSection;

    if (details.imgURL) {
      detailSection = (
        <div className="display-flex" style={{ alignItems: 'center' }}>
          <div
            className="si-thumbnail"
            style={{
              width: '56px',
              height: '42px',
              backgroundImage: `url(${details.imgURL})`,
            }}
            aria-labelledby="si-thumbnail--caption"
          />
          <small id="si-thumbnail--caption">{details.text}</small>
        </div>
      );
    } else {
      detailSection = <p className="si-details">{details.text}</p>;
    }

    return (
      <tr key={id} style={{ height: '80px' }}>
        <td style={{ paddingTop: '19px', verticalAlign: 'top' }}>{dateRequested}</td>
        <td style={{ paddingTop: '19px', verticalAlign: 'top' }}>{serviceItem}</td>
        <td style={{ paddingTop: '19px', verticalAlign: 'top' }}>{code}</td>
        <td style={{ verticalAlign: 'top' }}>{detailSection}</td>
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
    );
  });

  return (
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
        <tbody>{tableRows}</tbody>
      </table>
    </div>
  );
};

ServiceItemTableHasImg.propTypes = {
  serviceItems: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string,
      dateRequested: PropTypes.string,
      serviceItem: PropTypes.string,
      code: PropTypes.string,
      details: PropTypes.object,
    }),
  ).isRequired,
};

export default ServiceItemTableHasImg;
