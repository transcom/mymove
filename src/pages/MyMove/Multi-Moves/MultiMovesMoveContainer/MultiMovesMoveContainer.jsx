import React, { useState } from 'react';
import { useDispatch, connect } from 'react-redux';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';
import { useNavigate } from 'react-router';

import MultiMovesMoveInfoList from '../MultiMovesMoveInfoList/MultiMovesMoveInfoList';
import ButtonDropdownMenu from '../../../../components/ButtonDropdownMenu/ButtonDropdownMenu';

import styles from './MultiMovesMoveContainer.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import { customerRoutes } from 'constants/routes';
import { CHECK_SPECIAL_ORDERS_TYPES, SPECIAL_ORDERS_TYPES } from 'constants/orders';
import { setMoveId, setShowLoadingSpinner as setShowLoadingSpinnerAction } from 'store/general/actions';
import { ADVANCE_STATUSES } from 'constants/ppms';
import { onPacketDownloadSuccessHandler } from 'shared/AsyncPacketDownloadLink/AsyncPacketDownloadLink';
import { downloadPPMAOAPacket, downloadPPMPaymentPacket } from 'services/internalApi';
import { ppmShipmentStatuses } from 'constants/shipments';
import { setFlashMessage as setFlashMessageAction } from 'store/flash/actions';
import scrollToTop from 'shared/scrollToTop';
import { MOVE_STATUSES, SHIPMENT_TYPES } from 'shared/constants';
import { ButtonUsa as Button } from 'shared/standardUI/Buttons/ButtonUsa';

const MultiMovesMoveContainer = ({ moves, setFlashMessage, setShowLoadingSpinner }) => {
  const [expandedMoves, setExpandedMoves] = useState({});
  const navigate = useNavigate();
  const dispatch = useDispatch();

  // this expands the moves when the arrow is clicked
  const handleExpandClick = (index) => {
    setExpandedMoves((prev) => ({
      ...prev,
      [index]: !prev[index],
    }));
  };

  // handles the title of the shipment header below each move
  const generateShipmentTypeTitle = (shipmentType) => {
    if (shipmentType === SHIPMENT_TYPES.HHG) {
      return 'Household Goods';
    }
    if (shipmentType === SHIPMENT_TYPES.PPM) {
      return 'Personally Procured Move';
    }
    if (shipmentType === SHIPMENT_TYPES.NTS) {
      return 'Household Goods NTS';
    }
    if (shipmentType === SHIPMENT_TYPES.NTSR) {
      return 'Household Goods NTSR';
    }
    if (shipmentType === SHIPMENT_TYPES.MOBILE_HOME) {
      return 'Mobile Home';
    }
    if (shipmentType === SHIPMENT_TYPES.BOAT_HAUL_AWAY) {
      return 'Boat Haul Away';
    }
    if (shipmentType === SHIPMENT_TYPES.BOAT_TOW_AWAY) {
      return 'Boat Tow Away';
    }
    if (shipmentType === SHIPMENT_TYPES.UNACCOMPANIED_BAGGAGE) {
      return 'Unaccompanied Baggage';
    }
    return 'Shipment';
  };

  // sends user to the move page when clicking "Go to Move" btn
  const handleGoToMoveClick = (id) => {
    // When Go To Move is clicked store the moveId choosen in state
    dispatch(setMoveId(id));
    navigate(`${customerRoutes.MOVE_HOME_PAGE}/${id}`);
  };

  // this will determine what the PPM dropdown menu will show based on ppmShipment values present in the object
  const handlePPMDropdownOptions = (shipment) => {
    const { ppmShipment } = shipment;
    const dropdownOptions = {};

    if (
      ppmShipment?.advanceStatus === ADVANCE_STATUSES.APPROVED.apiValue &&
      ppmShipment?.status === ppmShipmentStatuses.CLOSEOUT_COMPLETE
    ) {
      dropdownOptions['AOA Packet'] = 'AOA Paperwork (PDF)';
      dropdownOptions['PPM Packet'] = 'PPM Packet';
    } else if (ppmShipment?.status === ppmShipmentStatuses.CLOSEOUT_COMPLETE) {
      dropdownOptions['PPM Packet'] = 'PPM Packet';
    } else {
      dropdownOptions['AOA Packet'] = 'AOA Paperwork (PDF)';
    }

    return Object.entries(dropdownOptions).map(([value], index) => ({
      id: index + 1,
      value,
    }));
  };

  // when an item is selected in the dropdown, this function handles API calls
  const handlePPMDropdownClick = (selectedItem, id) => {
    if (selectedItem.value === 'PPM Packet') {
      setShowLoadingSpinner(true, 'Downloading Payment Packet (PDF)...');
      downloadPPMPaymentPacket(id)
        .then((response) => {
          onPacketDownloadSuccessHandler(response);
          setFlashMessage('PPM_PACKET_DOWNLOAD_SUCCESS', 'success', 'PPM Packet successfully downloaded');
          setShowLoadingSpinner(false, null);
        })
        .catch(() => {
          setFlashMessage(
            'PPM_PACKET_DOWNLOAD_FAILURE',
            'error',
            'An error occurred when attempting download of PPM Packet',
          );
          setShowLoadingSpinner(false, null);
        });
    }
    if (selectedItem.value === 'AOA Packet') {
      setShowLoadingSpinner(true, 'Downloading AOA Paperwork (PDF)...');
      downloadPPMAOAPacket(id)
        .then((response) => {
          onPacketDownloadSuccessHandler(response);
          setFlashMessage('AOA_PACKET_DOWNLOAD_SUCCESS', 'success', 'AOA Packet successfully downloaded');
          setShowLoadingSpinner(false, null);
        })
        .catch(() => {
          setFlashMessage(
            'AOA_PACKET_DOWNLOAD_FAILURE',
            'error',
            'An error occurred when attempting download of AOA Packet',
          );
          setShowLoadingSpinner(false, null);
        });
    }
    scrollToTop();
  };

  const moveList = moves.map((m, index) => (
    <React.Fragment key={index}>
      <div className={styles.moveContainer}>
        <div className={styles.heading} key={index}>
          <h3>#{m.moveCode}</h3>
          {m.status === MOVE_STATUSES.CANCELED ? (
            <div className={styles.specialMoves}>Canceled</div>
          ) : (
            CHECK_SPECIAL_ORDERS_TYPES(m?.orders?.orders_type) && (
              <div className={styles.specialMoves}>{SPECIAL_ORDERS_TYPES[`${m?.orders?.orders_type}`]}</div>
            )
          )}
          <div className={styles.moveContainerButtons} data-testid="headerBtns">
            <Button
              data-testid="goToMoveBtn"
              className={styles.goToMoveBtn}
              onClick={() => {
                handleGoToMoveClick(m.id);
              }}
            >
              Go to Move
            </Button>
          </div>
          <FontAwesomeIcon
            className={styles.icon}
            icon={classnames({
              'chevron-up': expandedMoves[index],
              'chevron-down': !expandedMoves[index],
            })}
            data-testid="expand-icon"
            onClick={() => handleExpandClick(index)}
          />
        </div>
        <div className={styles.moveInfoList} data-testid="move-info-container">
          {expandedMoves[index] && (
            <div className={styles.moveInfoListExpanded}>
              <MultiMovesMoveInfoList move={m} />
              <h3 className={styles.shipmentH3}>Shipments</h3>
              {m.mtoShipments && m.mtoShipments.length > 0 ? (
                m.mtoShipments.map((s, sIndex) => (
                  <React.Fragment key={sIndex}>
                    <div className={styles.shipment} data-testid="shipment-container">
                      <ShipmentContainer
                        key={s.id}
                        shipmentType={s.shipmentType}
                        className={classnames(styles.previewShipment)}
                      >
                        <div className={styles.innerWrapper}>
                          <div className={styles.shipmentTypeHeading}>
                            <h4 data-testid="shipment-type">{generateShipmentTypeTitle(s.shipmentType)}</h4>
                            {s?.ppmShipment?.advanceStatus === ADVANCE_STATUSES.APPROVED.apiValue ||
                            s?.ppmShipment?.status === ppmShipmentStatuses.CLOSEOUT_COMPLETE ? (
                              <ButtonDropdownMenu
                                data-testid="ppmDownloadBtn"
                                title="Download"
                                items={handlePPMDropdownOptions(s)}
                                divClassName={styles.ppmDropdownBtn}
                                onItemClick={(e) => {
                                  handlePPMDropdownClick(e, s.ppmShipment.id);
                                }}
                                minimal
                              />
                            ) : null}
                            <h5>#{s.shipmentLocator}</h5>
                          </div>
                        </div>
                      </ShipmentContainer>
                    </div>
                  </React.Fragment>
                ))
              ) : (
                <div className={styles.shipment}>No shipments in move yet.</div>
              )}
            </div>
          )}
        </div>
      </div>
    </React.Fragment>
  ));

  return (
    <div data-testid="move-container" className={styles.movesContainer}>
      {moveList}
    </div>
  );
};

const mapDispatchToProps = {
  setFlashMessage: setFlashMessageAction,
  setShowLoadingSpinner: setShowLoadingSpinnerAction,
};

export default connect(() => ({}), mapDispatchToProps)(MultiMovesMoveContainer);
