import React, { useState, useEffect } from 'react';
import { Dropdown, Label } from '@trussworks/react-uswds';
import PropTypes from 'prop-types';
import { useNavigate, useParams, generatePath } from 'react-router';

import styles from './CreateShipmentServiceItemForm.module.scss';
import DestinationSITServiceItemForm from './DestinationSITServiceItemForm';
import OriginSITServiceItemForm from './OriginSITServiceItemForm';
import InternationalDestinationSITServiceItemForm from './InternationalDestinationSITServiceItemForm';
import InternationalOriginSITServiceItemForm from './InternationalOriginSITServiceItemForm';
import ShuttleSITServiceItemForm from './ShuttleSITServiceItemForm';
import DomesticCratingForm from './DomesticCratingForm';
import InternationalCratingForm from './InternationalCratingForm';
import InternationalShuttleServiceItemForm from './InternationalShuttleServiceItemForm';

import { isBooleanFlagEnabled } from 'utils/featureFlags';
import { ShipmentShape } from 'types/shipment';
import { createServiceItemModelTypes } from 'constants/prime';
import Shipment from 'components/PrimeUI/Shipment/Shipment';
import { FEATURE_FLAG_KEYS } from 'shared/constants';
import SectionWrapper from 'components/Shared/SectionWrapper/SectionWrapper';
import { primeSimulatorRoutes } from 'constants/routes';

const CreateShipmentServiceItemForm = ({ shipment, createServiceItemMutation }) => {
  const { moveCodeOrID } = useParams();
  const navigate = useNavigate();
  const handleCancel = () => {
    navigate(generatePath(primeSimulatorRoutes.VIEW_MOVE_PATH, { moveCodeOrID }));
  };
  const {
    MTOServiceItemOriginSIT,
    MTOServiceItemDestSIT,
    MTOServiceItemInternationalOriginSIT,
    MTOServiceItemInternationalDestSIT,
    MTOServiceItemDomesticShuttle,
    MTOServiceItemDomesticCrating,
    MTOServiceItemInternationalCrating,
    MTOServiceItemInternationalShuttle,
  } = createServiceItemModelTypes;
  const [selectedServiceItemType, setSelectedServiceItemType] = useState(MTOServiceItemOriginSIT);
  const [enableAlaskaFeatureFlag, setEnableAlaskaFeatureFlag] = useState(false);

  useEffect(() => {
    const fetchData = async () => {
      isBooleanFlagEnabled(FEATURE_FLAG_KEYS.ENABLE_ALASKA).then((res) => {
        setEnableAlaskaFeatureFlag(res);
      });
    };
    fetchData();
  }, []);

  const handleServiceItemTypeChange = (event) => {
    setSelectedServiceItemType(event.target.value);
  };

  return (
    <div className={styles.CreateShipmentServiceItemForm}>
      <Shipment shipment={shipment} />
      <SectionWrapper>
        <Label htmlFor="serviceItemType">Service item type</Label>
        <Dropdown id="serviceItemType" name="serviceItemType" onChange={handleServiceItemTypeChange}>
          <>
            <option value={MTOServiceItemOriginSIT}>Origin SIT</option>
            <option value={MTOServiceItemDestSIT}>Destination SIT</option>
            <option value={MTOServiceItemInternationalOriginSIT}>International Origin SIT</option>
            <option value={MTOServiceItemInternationalDestSIT}>International Destination SIT</option>
            <option value={MTOServiceItemDomesticShuttle}>Domestic Shuttle</option>
            <option value={MTOServiceItemInternationalShuttle}>International Shuttle</option>
            <option value={MTOServiceItemDomesticCrating}>Domestic Crating</option>
            {enableAlaskaFeatureFlag && (
              <option value={MTOServiceItemInternationalCrating}>International Crating</option>
            )}
          </>
        </Dropdown>
        {selectedServiceItemType === MTOServiceItemOriginSIT && (
          <OriginSITServiceItemForm
            shipment={shipment}
            submission={createServiceItemMutation}
            handleCancel={handleCancel}
          />
        )}
        {selectedServiceItemType === MTOServiceItemDestSIT && (
          <DestinationSITServiceItemForm
            shipment={shipment}
            submission={createServiceItemMutation}
            handleCancel={handleCancel}
          />
        )}

        {selectedServiceItemType === MTOServiceItemInternationalOriginSIT && (
          <InternationalOriginSITServiceItemForm
            shipment={shipment}
            submission={createServiceItemMutation}
            handleCancel={handleCancel}
          />
        )}
        {selectedServiceItemType === MTOServiceItemInternationalDestSIT && (
          <InternationalDestinationSITServiceItemForm
            shipment={shipment}
            submission={createServiceItemMutation}
            handleCancel={handleCancel}
          />
        )}

        {selectedServiceItemType === MTOServiceItemDomesticShuttle && (
          <ShuttleSITServiceItemForm
            shipment={shipment}
            submission={createServiceItemMutation}
            handleCancel={handleCancel}
          />
        )}
        {selectedServiceItemType === MTOServiceItemInternationalShuttle && (
          <InternationalShuttleServiceItemForm
            shipment={shipment}
            submission={createServiceItemMutation}
            handleCancel={handleCancel}
          />
        )}
        {selectedServiceItemType === MTOServiceItemDomesticCrating && (
          <DomesticCratingForm shipment={shipment} submission={createServiceItemMutation} handleCancel={handleCancel} />
        )}
        {enableAlaskaFeatureFlag && selectedServiceItemType === MTOServiceItemInternationalCrating && (
          <InternationalCratingForm
            shipment={shipment}
            submission={createServiceItemMutation}
            handleCancel={handleCancel}
          />
        )}
      </SectionWrapper>
    </div>
  );
};

CreateShipmentServiceItemForm.propTypes = {
  shipment: ShipmentShape.isRequired,
  createServiceItemMutation: PropTypes.func.isRequired,
};

export default CreateShipmentServiceItemForm;
