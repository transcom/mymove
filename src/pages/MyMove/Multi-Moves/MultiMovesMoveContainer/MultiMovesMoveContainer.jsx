import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';
import { useNavigate } from 'react-router';

import MultiMovesMoveInfoList from '../MultiMovesMoveInfoList/MultiMovesMoveInfoList';
import ButtonDropdownMenu from '../../../../components/ButtonDropdownMenu/ButtonDropdownMenu';

import styles from './MultiMovesMoveContainer.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';
import { customerRoutes } from 'constants/routes';
import { getMoveCodeLabel } from 'utils/shipmentDisplay';
import { SPECIAL_ORDERS_TYPES } from 'constants/orders';

const MultiMovesMoveContainer = ({ moves }) => {
  const [expandedMoves, setExpandedMoves] = useState({});
  const navigate = useNavigate();

  // this expands the moves when the arrow is clicked
  const handleExpandClick = (index) => {
    setExpandedMoves((prev) => ({
      ...prev,
      [index]: !prev[index],
    }));
  };

  // when an item is selected in the dropdown, this function will handle that logic
  const handleDropdownItemClick = (selectedItem) => {
    return selectedItem.value;
  };

  const dropdownMenuItems = [
    {
      id: 1,
      value: 'PCS Orders',
    },
    {
      id: 2,
      value: 'PPM Packet',
    },
  ];

  // handles the title of the shipment header below each move
  const generateShipmentTypeTitle = (shipmentType) => {
    if (shipmentType === 'HHG') {
      return 'Household Goods';
    }
    if (shipmentType === 'PPM') {
      return 'Personally Procured Move';
    }
    if (shipmentType === 'HHG_INTO_NTS_DOMESTIC') {
      return 'Household Goods NTS';
    }
    if (shipmentType === 'HHG_OUTOF_NTS_DOMESTIC') {
      return 'Household Goods NTSR';
    }
    if (shipmentType === 'MOTORHOME') {
      return 'Motorhome';
    }
    if (shipmentType === 'BOAT_HAUL_AWAY') {
      return 'Boat Haul Away';
    }
    if (shipmentType === 'BOAT_TOW_AWAY') {
      return 'Boat Tow Away';
    }
    return 'Shipment';
  };

  // sends user to the move page when clicking "Go to Move" btn
  const handleGoToMoveClick = (id) => {
    navigate(`${customerRoutes.MOVE_HOME_PAGE}/${id}`);
  };

  const moveList = moves.map((m, index) => (
    <React.Fragment key={index}>
      <div className={styles.moveContainer}>
        <div className={styles.heading} key={index}>
          <h3>#{m.moveCode}</h3>
          {['BLUEBARK', 'WOUNDED_WARRIOR'].includes(m?.orders?.orders_type) ? (
            <div className={styles.specialMoves}>{SPECIAL_ORDERS_TYPES[`${m?.orders?.orders_type}`]}</div>
          ) : null}
          <div className={styles.moveContainerButtons} data-testid="headerBtns">
            <ButtonDropdownMenu
              data-testid="downloadBtn"
              title="Download"
              items={dropdownMenuItems}
              divClassName={styles.dropdownBtn}
              onItemClick={handleDropdownItemClick}
              outline
            />
            <Button
              data-testid="goToMoveBtn"
              className={styles.goToMoveBtn}
              secondary
              outline
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
                    <div className={styles.shipment}>
                      <ShipmentContainer
                        key={s.id}
                        shipmentType={s.shipmentType}
                        className={classnames(styles.previewShipment)}
                      >
                        <div className={styles.innerWrapper}>
                          <div className={styles.shipmentTypeHeading}>
                            <h4>{generateShipmentTypeTitle(s.shipmentType)}</h4>
                            <h5>#{getMoveCodeLabel(s.id)}</h5>
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

export default MultiMovesMoveContainer;
