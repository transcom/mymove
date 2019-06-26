import React from 'react';
import { string, number, bool } from 'prop-types';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import carImg from 'shared/images/car_mobile.png';
import boxTruckImg from 'shared/images/box_truck_mobile.png';
import carTrailerImg from 'shared/images/car-trailer_mobile.png';
import deleteButtonImg from 'shared/images/delete-doc-button.png';
import { intToOrdinal } from '../utility';

const WEIGHT_TICKET_IMAGES = {
  CAR: carImg,
  BOX_TRUCK: boxTruckImg,
  CAR_TRAILER: carTrailerImg,
};

const MissingLabel = ({ children }) => (
  <p className="missing-label">
    <em>{children}</em>
  </p>
);

const WeightTicketListItem = ({
  vehicle_options,
  vehicle_nickname,
  num,
  empty_weight,
  full_weight,
  empty_weight_ticket_missing,
  full_weight_ticket_missing,
  trailer_ownership_missing,
}) => (
  <div className="ticket-item" style={{ display: 'flex' }}>
    {/* size of largest of the images */}
    <div style={{ minWidth: 95 }}>
      {/*eslint-disable security/detect-object-injection*/}
      <img className="weight-ticket-image" src={WEIGHT_TICKET_IMAGES[vehicle_options]} alt={vehicle_options} />
    </div>
    <div style={{ flex: 1 }}>
      <div className="weight-li-item-container">
        <h4>
          {vehicle_nickname} ({intToOrdinal(num + 1)} set)
        </h4>
        <img alt="delete document button" onClick={() => console.log('lol')} src={deleteButtonImg} />
      </div>
      {empty_weight_ticket_missing ? (
        <MissingLabel>
          Missing empty weight ticket{' '}
          <FontAwesomeIcon style={{ color: 'red' }} className="icon" icon={faExclamationCircle} />
        </MissingLabel>
      ) : (
        <p>Empty weight ticket {empty_weight} lbs</p>
      )}
      {full_weight_ticket_missing ? (
        <MissingLabel>
          Missing full weight ticket{' '}
          <FontAwesomeIcon style={{ color: 'red' }} className="icon" icon={faExclamationCircle} />
        </MissingLabel>
      ) : (
        <p>Full weight ticket {full_weight} lbs</p>
      )}
      {vehicle_options === 'CAR_TRAILER' &&
        trailer_ownership_missing && (
          <MissingLabel>
            Missing ownership documentation{' '}
            <FontAwesomeIcon style={{ color: 'red' }} className="icon" icon={faExclamationCircle} />
          </MissingLabel>
        )}
      {vehicle_options === 'CAR_TRAILER' && !trailer_ownership_missing && <p>Ownership documentation</p>}
    </div>
  </div>
);

WeightTicketListItem.propTypes = {
  vehicle_options: string.isRequired,
  vehicle_nickname: string.isRequired,
  num: number.isRequired,
  empty_weight: number.isRequired,
  full_weight: number.isRequired,
  empty_weight_ticket_missing: bool.isRequired,
  full_weight_ticket_missing: bool.isRequired,
  trailer_ownership_missing: bool.isRequired,
};
export default WeightTicketListItem;
