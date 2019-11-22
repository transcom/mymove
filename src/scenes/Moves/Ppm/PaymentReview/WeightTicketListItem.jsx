import React, { Component } from 'react';
import { forEach } from 'lodash';
import { string, number, bool } from 'prop-types';
import FontAwesomeIcon from '@fortawesome/react-fontawesome';
import faExclamationCircle from '@fortawesome/fontawesome-free-solid/faExclamationCircle';
import carImg from 'shared/images/car_mobile.png';
import boxTruckImg from 'shared/images/box_truck_mobile.png';
import carTrailerImg from 'shared/images/car-trailer_mobile.png';
import { formatToOrdinal } from 'shared/formatters';
import deleteButtonImg from 'shared/images/delete-doc-button.png';
import AlertWithDeleteConfirmation from 'shared/AlertWithDeleteConfirmation';
import { UPLOAD_SCAN_STATUS } from 'shared/constants';

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

class WeightTicketListItem extends Component {
  state = {
    showDeleteConfirmation: false,
  };

  areUploadsInfected = uploads => {
    let isInfected = false;
    forEach(uploads, function(upload) {
      if (upload.status === UPLOAD_SCAN_STATUS.INFECTED) {
        isInfected = true;
      }
    });
    return isInfected;
  };

  toggleShowConfirmation = () => {
    const { showDeleteConfirmation } = this.state;
    this.setState({ showDeleteConfirmation: !showDeleteConfirmation });
  };

  render() {
    const {
      id,
      empty_weight_ticket_missing,
      empty_weight,
      full_weight_ticket_missing,
      full_weight,
      num,
      trailer_ownership_missing,
      vehicle_nickname,
      vehicle_options,
      showDelete,
      deleteDocumentListItem,
      isWeightTicketSet,
      uploads,
    } = this.props;
    const { showDeleteConfirmation } = this.state;
    const isInfected = this.areUploadsInfected(uploads);
    return (
      <div className="ticket-item" style={{ display: 'flex' }}>
        {/* size of largest of the images */}
        <div style={{ minWidth: 95 }}>
          {/*eslint-disable security/detect-object-injection*/}
          <img className="weight-ticket-image" src={WEIGHT_TICKET_IMAGES[vehicle_options]} alt={vehicle_options} />
        </div>
        <div style={{ flex: 1 }}>
          <div className="weight-li-item-container">
            <h4>
              {isWeightTicketSet && (
                <>
                  {vehicle_nickname}
                  {formatToOrdinal(num + 1)} set
                </>
              )}
            </h4>
            {showDelete && (
              <img
                alt="delete document button"
                data-cy="delete-ticket"
                onClick={this.toggleShowConfirmation}
                src={deleteButtonImg}
              />
            )}
          </div>
          {isInfected && (
            <>
              <div className="infected-indicator">
                <strong>Delete this file, take a photo of the document, then upload that</strong>
              </div>
            </>
          )}
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
          {vehicle_options === 'CAR_TRAILER' && trailer_ownership_missing && (
            <MissingLabel>
              Missing ownership documentation{' '}
              <FontAwesomeIcon style={{ color: 'red' }} className="icon" icon={faExclamationCircle} />
            </MissingLabel>
          )}
          {vehicle_options === 'CAR_TRAILER' && !trailer_ownership_missing && <p>Ownership documentation</p>}
          {showDeleteConfirmation && (
            <AlertWithDeleteConfirmation
              heading="Delete this document?"
              message="This action cannot be undone."
              deleteActionHandler={() => deleteDocumentListItem(id)}
              cancelActionHandler={this.toggleShowConfirmation}
              type="weight-ticket-list-alert"
            />
          )}
        </div>
      </div>
    );
  }
}

WeightTicketListItem.propTypes = {
  id: string.isRequired,
  num: number.isRequired,
  isWeightTicketSet: bool.isRequired,
};

WeightTicketListItem.defaultProps = {
  showDelete: false,
};
export default WeightTicketListItem;
