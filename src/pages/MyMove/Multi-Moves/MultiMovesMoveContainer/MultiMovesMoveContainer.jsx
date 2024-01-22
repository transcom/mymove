import React, { useState } from 'react';
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome';
import classnames from 'classnames';
import { Button } from '@trussworks/react-uswds';

import MultiMovesMoveInfoList from '../MultiMovesMoveInfoList/MultiMovesMoveInfoList';
import ButtonDropdownMenu from '../ButtonDropdownMenu/ButtonDropdownMenu';

import styles from './MultiMovesMoveContainer.module.scss';

import ShipmentContainer from 'components/Office/ShipmentContainer/ShipmentContainer';

const MultiMovesMoveContainer = ({ moves }) => {
  const [expandedMoves, setExpandedMoves] = useState({});

  const handleExpandClick = (index) => {
    setExpandedMoves((prev) => ({
      ...prev,
      [index]: !prev[index],
    }));
  };

  const handleDropdownItemClick = (selectedItem) => {
    console.log(`${selectedItem.value}`);
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

  const moveList = moves.map((m, index) => (
    <React.Fragment key={index}>
      <div className={styles.moveContainer}>
        <div className={styles.heading} key={index}>
          <h3>#{m.moveCode}</h3>
          {m.status !== 'COMPLETED' ? (
            <Button className={styles.goToMoveBtn} secondary outline>
              Go to Move
            </Button>
          ) : (
            <ButtonDropdownMenu
              title="Download"
              items={dropdownMenuItems}
              divClassName={styles.dropdownBtn}
              onItemClick={handleDropdownItemClick}
              outline
            />
          )}
          <FontAwesomeIcon
            className={styles.icon}
            icon={classnames({
              'chevron-up': expandedMoves[index],
              'chevron-down': !expandedMoves[index],
            })}
            onClick={() => handleExpandClick(index)}
          />
        </div>
        <div className={styles.moveInfoList} data-testid="move-info-container">
          {expandedMoves[index] && (
            <div>
              <MultiMovesMoveInfoList move={m} />
              <h3 className={styles.shipmentH3}>Shipments</h3>
              {m.mtoShipments.map((s, sIndex) => (
                <React.Fragment key={sIndex}>
                  <div className={styles.shipment}>
                    <ShipmentContainer
                      key={s.id}
                      shipmentType={s.shipmentType}
                      className={classnames(styles.previewShipment)}
                    >
                      <div className={styles.innerWrapper}>
                        <div className={styles.shipmentTypeHeading}>
                          <h4>{s.shipmentType}</h4>
                          <h5>#{m.moveCode}</h5>
                        </div>
                      </div>
                    </ShipmentContainer>
                  </div>
                </React.Fragment>
              ))}
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
