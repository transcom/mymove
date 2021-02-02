import React from 'react';
import { action } from '@storybook/addon-actions';

import { StorageInfoModal } from './StorageInfoModal/StorageInfoModal';
import { MoveInfoModal } from './MoveInfoModal/MoveInfoModal';
import { EulaModal } from './EulaModal/EulaModal';

export default {
  title: 'Customer Components/Modals',
};

const eulaProps = {
  acceptTerms: action('clicked'),
};

export const StorageInfoModalStory = () => <StorageInfoModal />;
export const MoveInfoModalStory = () => <MoveInfoModal />;
// eslint-disable-next-line react/jsx-props-no-spreading
export const EulaModalStory = () => <EulaModal {...eulaProps} />;
