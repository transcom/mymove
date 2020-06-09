import React, { useState } from 'react';
import * as PropTypes from 'prop-types';
import classNames from 'classnames/bind';
import { Button, Checkbox, Fieldset, Modal, Overlay, ModalContainer } from '@trussworks/react-uswds';

import styles from './requestedShipments.module.scss';

const cx = classNames.bind(styles);

const RequestedShipments = ({ children }) => {
  const [isModalVisible, setIsModalVisible] = useState(false);
  const [isApproveButtonDisabled, setIsApproveButtonDisabled] = useState(true);
  const handleApprovalClick = () => {
    setIsModalVisible(true);
  };

  const onChange = () => {
    const shipmentManagementFeeChecked = document.getElementById('shipmentManagementFee').valueOf().checked;
    const counselingFeeChecked = document.getElementById('counselingFee').valueOf().checked;
    const hhgShipmentChecked = document.getElementById('shipment-display-checkbox-hhg').valueOf().checked;

    if (hhgShipmentChecked && (counselingFeeChecked || shipmentManagementFeeChecked)) {
      setIsApproveButtonDisabled(false);
    } else {
      setIsApproveButtonDisabled(true);
    }
  };

  return (
    <div className={`${cx('requested-shipments')} container`} data-cy="requested-shipments">
      <div id="approvalConfirmationModal" style={{ display: isModalVisible ? 'block' : 'none' }}>
        <Overlay />
        <ModalContainer>
          <Modal>
            <h1>Preview and post move task order</h1>
            <button type="button" onClick={() => setIsModalVisible(false)}>
              Temporary Close Button
            </button>
          </Modal>
        </ModalContainer>
      </div>
      <h4>Requested shipments</h4>
      <div className={`${cx('__content')}`}>
        {children &&
          React.Children.map(children, (child, index) => (
            // eslint-disable-next-line react/no-array-index-key
            <div key={index} className={`${cx('__item')}`}>
              {React.cloneElement(child, { onChange })}
            </div>
          ))}
      </div>
      <div>
        <h3>Add service items to this move</h3>
        <span>{isModalVisible}</span>
        <Fieldset legend="MTO service items" legendSrOnly id="input-type-fieldset">
          <Checkbox
            id="shipmentManagementFee"
            label="Shipment management fee"
            name="shipmentManagementFee"
            value="true"
            onChange={onChange}
          />
          <Checkbox id="counselingFee" label="Counseling fee" name="counselingFee" value="true" onChange={onChange} />
        </Fieldset>
        <Button
          id="shipmentApproveButton"
          className={`${cx('usa-button--small')} usa-button--icon`}
          onClick={handleApprovalClick}
          type="button"
          disabled={isApproveButtonDisabled}
        >
          <span>Approve selected shipments</span>
        </Button>
      </div>
    </div>
  );
};

RequestedShipments.propTypes = {
  children: PropTypes.oneOfType([PropTypes.element, PropTypes.arrayOf(PropTypes.element)]),
};

export default RequestedShipments;
