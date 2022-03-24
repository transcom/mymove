import React from 'react';

import { StorageInfoModal } from './StorageInfoModal/StorageInfoModal';
import { MoveInfoModal } from './MoveInfoModal/MoveInfoModal';

import { AddShipmentModal } from 'components/Customer/Review/AddShipmentModal/AddShipmentModal';

export default {
  title: 'Customer Components/Modals',
};

export const StorageInfoModalStory = () => <StorageInfoModal />;
export const MoveInfoModalStory = () => <MoveInfoModal />;
export const AddShipmentModalStory = () => <AddShipmentModal />;
