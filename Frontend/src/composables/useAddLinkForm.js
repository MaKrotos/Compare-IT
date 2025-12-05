import { ref, computed } from "vue";
import { getCachedPreviewData } from "@/utils/preview.js";

/**
 * Композабл для управления логикой формы добавления ссылок
 */
export function useAddLinkForm() {
  const links = ref([
    { url: "", price: "", imageUrls: "", title: "", description: "" },
  ]);
  const isDragOver = ref(false);
  const isLoading = ref(false);

  /**
   * Обновляет ссылку по индексу
   * @param {number} index - Индекс ссылки
   * @param {object} newLink - Новые данные ссылки
   */
  const updateLink = (index, newLink) => {
    console.log(
      "updateLink вызван с индексом:",
      index,
      "и новыми данными:",
      newLink
    );
    links.value[index] = { ...links.value[index], ...newLink };
    console.log("Обновленные ссылки:", links.value);
  };

  /**
   * Добавляет новое поле для ввода ссылки
   */
  const addLinkField = () => {
    console.log("addLinkField вызван");
    if (links.value.length < 10) {
      links.value.push({
        url: "",
        price: "",
        imageUrls: "",
        title: "",
        description: "",
      });
      console.log(
        "Добавлено новое поле для ссылки. Текущие ссылки:",
        links.value
      );
    } else {
      console.log("Достигнуто максимальное количество полей для ссылок (10)");
    }
  };

  /**
   * Удаляет поле по индексу
   * @param {number} index - Индекс поля для удаления
   */
  const removeLink = (index) => {
    console.log("removeLink вызван с индексом:", index);
    const newLinks = links.value.filter((_, i) => i !== index);
    links.value = newLinks.length
      ? newLinks
      : [{ url: "", price: "", imageUrls: "", title: "", description: "" }];
    console.log("Ссылки после удаления:", links.value);
  };

  /**
   * Обработчик добавления товаров
   * @param {function} onAddLinks - Функция обратного вызова для добавления ссылок
   */
  const handleAddItems = (onAddLinks) => {
    console.log("Ссылки перед фильтрацией:", links.value);
    const validLinks = links.value.filter((link) => link.url.trim());
    console.log("Валидные ссылки:", validLinks);

    if (validLinks.length > 0) {
      // Преобразуем imageUrls в массив
      const processedLinks = validLinks.map((link) => {
        let images = [];
        if (link.imageUrls) {
          images = link.imageUrls
            .split("\n")
            .map((url) => url.trim())
            .filter((url) => url);
        }
        return {
          ...link,
          images: images,
          // Используем заголовок из формы или извлекаем из URL
          title:
            link.title && link.title !== "Загрузка..."
              ? link.title
              : "Без названия",
        };
      });

      console.log("Добавляем товары:", processedLinks);
      onAddLinks(processedLinks);
      // Очищаем форму после добавления
      links.value = [
        { url: "", price: "", imageUrls: "", title: "", description: "" },
      ];
    } else {
      console.log("Нет валидных ссылок для добавления");
    }
  };

  /**
   * Обработка вставки изображений через Ctrl+V
   * @param {number} index - Индекс ссылки
   * @param {ClipboardEvent} event - Событие вставки
   */
  const handleImagePaste = (index, event) => {
    console.log("handleImagePaste вызван с индексом:", index);
    const items = (event.clipboardData || event.originalEvent.clipboardData)
      .items;

    for (let i = 0; i < items.length; i++) {
      const item = items[i];

      // Если вставлен файл изображения
      if (item.type.indexOf("image") !== -1) {
        const blob = item.getAsFile();
        const reader = new FileReader();

        reader.onload = (e) => {
          // Создаем data URL изображения
          const dataURL = e.target.result;
          console.log("Изображение вставлено как data URL:", dataURL);

          // Добавляем к существующим ссылкам
          const currentUrls = links.value[index].imageUrls || "";
          const separator = currentUrls ? "\n" : "";
          const newUrls = currentUrls + separator + dataURL;
          updateLink(index, { imageUrls: newUrls });
        };

        reader.readAsDataURL(blob);
      }

      // Если вставлена ссылка
      if (item.type === "text/plain") {
        item.getAsString((text) => {
          console.log("Текст вставлен:", text);
          // Проверяем, является ли текст URL
          if (text.match(/^https?:\/\/.+\.(jpg|jpeg|png|gif|webp)/i)) {
            console.log("Вставлен URL изображения:", text);
            // Добавляем к существующим ссылкам
            const currentUrls = links.value[index].imageUrls || "";
            const separator = currentUrls ? "\n" : "";
            const newUrls = currentUrls + separator + text;
            updateLink(index, { imageUrls: newUrls });
          }
        });
      }
    }
  };

  // Обработка событий drag
  const handleDragEnter = (event) => {
    console.log("handleDragEnter вызван");
    event.preventDefault();
    isDragOver.value = true;
  };

  const handleDragLeave = (event) => {
    console.log("handleDragLeave вызван");
    event.preventDefault();
    isDragOver.value = false;
  };

  // Обработка drag and drop
  const handleDrop = (event, index) => {
    console.log("handleDrop вызван с индексом:", index);
    event.preventDefault();
    isDragOver.value = false;

    const files = event.dataTransfer.files;
    if (files && files.length > 0) {
      console.log("Файлы перетянуты:", files);
      for (let i = 0; i < files.length; i++) {
        const file = files[i];
        if (file.type.startsWith("image/")) {
          const reader = new FileReader();
          reader.onload = (e) => {
            // Создаем data URL изображения
            const dataURL = e.target.result;
            console.log("Изображение перетянуто как data URL:", dataURL);

            // Добавляем к существующим ссылкам
            const currentUrls = links.value[index].imageUrls || "";
            const separator = currentUrls ? "\n" : "";
            const newUrls = currentUrls + separator + dataURL;
            updateLink(index, { imageUrls: newUrls });
          };
          reader.readAsDataURL(file);
        }
      }
    }
  };

  // Вычисляемое свойство для проверки валидности формы
  const isFormValid = computed(() => {
    const isValid = links.value.some((l) => l.url.trim());
    console.log("isFormValid вычисляется. Результат:", isValid);
    return isValid;
  });

  return {
    links,
    isDragOver,
    isLoading,
    updateLink,
    addLinkField,
    removeLink,
    handleAddItems,
    handleImagePaste,
    handleDragEnter,
    handleDragLeave,
    handleDrop,
    isFormValid,
  };
}
