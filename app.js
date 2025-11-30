// Конфигурация API
// Если фронтенд и бэкенд на одном сервере, используем относительный путь
const API_URL = '/api';

// Получение ссылок на элементы DOM
const createPostForm = document.getElementById('createPostForm');
const postsContainer = document.getElementById('postsContainer');
const editModal = document.getElementById('editModal');
const editPostForm = document.getElementById('editPostForm');
const closeModal = document.querySelector('.close');

// Инициализация приложения при загрузке страницы
document.addEventListener('DOMContentLoaded', () => {
    loadPosts();
    setupEventListeners();
});

// Настройка обработчиков событий
function setupEventListeners() {
    // Обработчик отправки формы создания поста
    createPostForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        await createPost();
    });

    // Обработчик отправки формы редактирования поста
    editPostForm.addEventListener('submit', async (e) => {
        e.preventDefault();
        await updatePost();
    });

    // Закрытие модального окна при клике на крестик
    closeModal.addEventListener('click', () => {
        editModal.style.display = 'none';
    });

    // Закрытие модального окна при клике вне его области
    window.addEventListener('click', (e) => {
        if (e.target === editModal) {
            editModal.style.display = 'none';
        }
    });
}

// Загрузка всех постов с сервера
async function loadPosts() {
    try {
        // Показываем индикатор загрузки
        postsContainer.innerHTML = '<div class="loading">Загрузка постов...</div>';

        // Выполняем GET-запрос к API для получения всех постов
        const response = await fetch(`${API_URL}/posts`);

        if (!response.ok) {
            throw new Error('Ошибка при загрузке постов');
        }

        const posts = await response.json();

        // Если постов нет, показываем пустое состояние
        if (!posts || posts.length === 0) {
            postsContainer.innerHTML = `
                <div class="empty-state">
                    <h3>Пока нет постов</h3>
                    <p>Создайте первый пост, используя форму выше!</p>
                </div>
            `;
            return;
        }

        // Отображаем посты на странице
        displayPosts(posts);
    } catch (error) {
        console.error('Ошибка:', error);
        postsContainer.innerHTML = `
            <div class="message error">
                Не удалось загрузить посты. Проверьте, что бэкенд запущен на ${API_URL}
            </div>
        `;
    }
}

// Отображение постов в DOM
function displayPosts(posts) {
    // Сортируем посты по дате создания (новые сверху)
    posts.sort((a, b) => new Date(b.created_at) - new Date(a.created_at));

    // Генерируем HTML для каждого поста и объединяем в одну строку
    postsContainer.innerHTML = posts.map(post => `
        <div class="post-card" data-id="${post.id}">
            <h3>${escapeHtml(post.title)}</h3>
            <div class="post-meta">
                Автор: ${escapeHtml(post.author)} | 
                Дата: ${formatDate(post.created_at)}
            </div>
            <div class="post-content">${escapeHtml(post.content)}</div>
            <div class="post-actions">
                <button class="btn btn-secondary" onclick="openEditModal(${post.id})">
                    Редактировать
                </button>
                <button class="btn btn-danger" onclick="deletePost(${post.id})">
                    Удалить
                </button>
            </div>
        </div>
    `).join('');
}

// Создание нового поста
async function createPost() {
    // Получаем данные из формы
    const title = document.getElementById('postTitle').value;
    const author = document.getElementById('postAuthor').value;
    const content = document.getElementById('postContent').value;

    try {
        // Отправляем POST-запрос на сервер для создания нового поста
        const response = await fetch(`${API_URL}/posts`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ title, author, content })
        });

        if (!response.ok) {
            throw new Error('Ошибка при создании поста');
        }

        // Очищаем форму после успешного создания
        createPostForm.reset();

        // Показываем сообщение об успехе
        showMessage('Пост успешно создан!', 'success');

        // Перезагружаем список постов
        await loadPosts();
    } catch (error) {
        console.error('Ошибка:', error);
        showMessage('Не удалось создать пост', 'error');
    }
}

// Открытие модального окна для редактирования поста
async function openEditModal(postId) {
    try {
        // Загружаем данные поста с сервера
        const response = await fetch(`${API_URL}/posts/${postId}`);

        if (!response.ok) {
            throw new Error('Ошибка при загрузке поста');
        }

        const post = await response.json();

        // Заполняем форму редактирования данными поста
        document.getElementById('editPostId').value = post.id;
        document.getElementById('editPostTitle').value = post.title;
        document.getElementById('editPostAuthor').value = post.author;
        document.getElementById('editPostContent').value = post.content;

        // Показываем модальное окно
        editModal.style.display = 'block';
    } catch (error) {
        console.error('Ошибка:', error);
        showMessage('Не удалось загрузить пост для редактирования', 'error');
    }
}

// Обновление существующего поста
async function updatePost() {
    const postId = document.getElementById('editPostId').value;
    const title = document.getElementById('editPostTitle').value;
    const author = document.getElementById('editPostAuthor').value;
    const content = document.getElementById('editPostContent').value;

    try {
        // Отправляем PUT-запрос на сервер для обновления поста
        const response = await fetch(`${API_URL}/posts/${postId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ title, author, content })
        });

        if (!response.ok) {
            throw new Error('Ошибка при обновлении поста');
        }

        // Закрываем модальное окно
        editModal.style.display = 'none';

        // Показываем сообщение об успехе
        showMessage('Пост успешно обновлен!', 'success');

        // Перезагружаем список постов
        await loadPosts();
    } catch (error) {
        console.error('Ошибка:', error);
        showMessage('Не удалось обновить пост', 'error');
    }
}

// Удаление поста
async function deletePost(postId) {
    // Запрашиваем подтверждение у пользователя
    if (!confirm('Вы уверены, что хотите удалить этот пост?')) {
        return;
    }

    try {
        // Отправляем DELETE-запрос на сервер
        const response = await fetch(`${API_URL}/posts/${postId}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            throw new Error('Ошибка при удалении поста');
        }

        // Показываем сообщение об успехе
        showMessage('Пост успешно удален!', 'success');

        // Перезагружаем список постов
        await loadPosts();
    } catch (error) {
        console.error('Ошибка:', error);
        showMessage('Не удалось удалить пост', 'error');
    }
}

// Вспомогательные функции

// Функция для отображения сообщений пользователю
function showMessage(text, type) {
    const messageDiv = document.createElement('div');
    messageDiv.className = `message ${type}`;
    messageDiv.textContent = text;

    // Вставляем сообщение в начало контейнера
    const container = document.querySelector('.container');
    container.insertBefore(messageDiv, container.firstChild);

    // Автоматически удаляем сообщение через 3 секунды
    setTimeout(() => {
        messageDiv.remove();
    }, 3000);
}

// Форматирование даты в читаемый вид
function formatDate(dateString) {
    const date = new Date(dateString);
    const options = {
        year: 'numeric',
        month: 'long',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
    };
    return date.toLocaleDateString('ru-RU', options);
}

// Защита от XSS-атак - экранирование HTML-символов
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}