import { ModalName } from '../page';

interface ImportWallet {
  onSwitchModal(modalName: ModalName): void;
}

const ImportWallet = ({}: ImportWallet) => {
  return (
    <div>
      <h1>This is form import wallet</h1>
    </div>
  );
};

export default ImportWallet;
